package service

import (
	"bufio"
	"context"
	"errors"
	"io"
	"log/slog"
	"os/exec"
	"sync"

	"github.com/mirhijinam/pg-project/internal/model"
	"github.com/mirhijinam/pg-project/internal/repository"
)

type CommandService struct {
	repo    Repo
	Workers sync.Map
	Pool    *Pool
}

type Repo interface {
	Create(context.Context, *model.Command) error
	SetError(context.Context, int, string) error
	SetSuccess(context.Context, int) error
	GetList(context.Context) ([]model.Command, error)
	GetOne(context.Context, int) (model.Command, error)
	Writer(context.Context, int) repository.WriterFunc
}

func New(repo Repo, pool *Pool) *CommandService {
	return &CommandService{
		repo:    repo,
		Pool:    pool,
		Workers: sync.Map{},
	}
}

func (cs *CommandService) DeleteCommand(id int) error {
	defer func() {
		cs.Workers.Delete(id)
	}()

	loadedExecCmd, ok := cs.Workers.Load(id)
	if !ok {
		slog.Error("failed to find cmd with such id")
		return errors.New("failed to find cmd with such id")
	}

	err := loadedExecCmd.(*model.CommandExec).Exec.Process.Kill()

	return err

}

func (cs *CommandService) CreateCommand(cmd *model.Command, isLong bool) error {
	slog.Info("in the service.CreateCommand")
	ctx := context.Background()

	err := cs.repo.Create(ctx, cmd)
	if err != nil {
		return err
	}

	cs.Pool.Go(func() {
		slog.Info("in the func of the Pool.Go")
		err := cs.execCommand(ctx, *cmd, isLong)
		if err != nil {
			slog.ErrorContext(ctx, "failed to exec the command", err)
		}
	})

	return nil
}

func (cs *CommandService) execCommand(ctx context.Context, cmd model.Command, isLong bool) (err error) {
	var wg sync.WaitGroup
	errch := make(chan error, 1)

	defer func() {
		wg.Wait()
		close(errch)
		receivedErr := <-errch
		if receivedErr != nil {
			err = errors.Join(err, cs.repo.SetError(ctx, cmd.Id, receivedErr.Error()))
			slog.Error("failed to execute command", "error", receivedErr.Error())
		} else {
			if err = cs.repo.SetSuccess(ctx, cmd.Id); err != nil {
				slog.Error("failed to set command success in db", "error", err.Error())
			}
		}
	}()

	slog.Info("in the func execCommand")
	execCmd := exec.Command("/bin/sh", "-c", cmd.Raw)

	stdout, err := execCmd.StdoutPipe()
	if err != nil {
		slog.Error("failed to create stdout pipe", "error", err.Error())
		return err
	}

	err = execCmd.Start()
	if err != nil {
		slog.Error("failed to start command", "error", err.Error())
		return err
	}

	cs.Workers.Store(cmd.Id, &model.CommandExec{
		Cmd:  cmd,
		Exec: execCmd,
	})

	scanner := bufio.NewScanner(stdout)
	wg.Add(1)
	if isLong {
		go cs.long(ctx, stdout, scanner, execCmd, errch, &wg, &cmd)
	} else {
		go cs.short(ctx, stdout, scanner, execCmd, errch, &wg, &cmd)
	}

	return nil
}

func (cs *CommandService) short(ctx context.Context, stdout io.ReadCloser, scanner *bufio.Scanner, execCmd *exec.Cmd, errch chan error, wg *sync.WaitGroup, cmd *model.Command) {
	slog.Info("in the func short")
	defer stdout.Close()
	defer wg.Done()

	var err error
	defer func() {
		errch <- err
	}()

	var accumulate string
	for scanner.Scan() {
		line := scanner.Text()
		accumulate += line + "\n"
	}

	err = execCmd.Wait()
	if err == nil {
		slog.Info("successfully executed and wrote into db", slog.String("short cmd", cmd.Raw))
	} else {
		slog.Error("command completed with error", "error", err.Error())
		return
	}

	_, err = cs.repo.Writer(ctx, cmd.Id)([]byte(accumulate))
	if err != nil {
		slog.Error("failed to write cmd result into db", "error", err.Error())
		return
	}
}

func (cs *CommandService) long(ctx context.Context, stdout io.ReadCloser, scanner *bufio.Scanner, execCmd *exec.Cmd, errch chan error, wg *sync.WaitGroup, cmd *model.Command) {
	slog.Info("in the func long")
	defer stdout.Close()
	defer wg.Done()

	var err error
	defer func() {
		errch <- err
	}()

	for scanner.Scan() {
		line := scanner.Text()
		_, err := cs.repo.Writer(ctx, cmd.Id)([]byte(line + "\n"))
		if err != nil {
			slog.Error("failed to write cmd result into db", "error", err.Error())
			return
		}
	}

	err = execCmd.Wait()
	if err == nil {
		slog.Info("successfully executed and wrote into db", slog.String("long cmd", cmd.Raw))
	} else {
		slog.Error("failed to set command success in db", "error", err.Error())
		return
	}
}

func (cs *CommandService) GetCommand(ctx context.Context, id int) (model.Command, error) {
	return cs.repo.GetOne(ctx, id)
}

func (cs *CommandService) GetCommandList(ctx context.Context) ([]model.Command, error) {
	return cs.repo.GetList(ctx)
}
