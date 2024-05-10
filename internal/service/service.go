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
	Repo    Repo
	Workers sync.Map
	Pool    *Pool
}

type Repo interface {
	Create(context.Context, *model.Command) error
	SetError(context.Context, int, string) error
	SetSuccess(context.Context, int) error
	List(context.Context) ([]model.Command, error)
	Get(context.Context, string) (model.Command, error)
	Writer(context.Context, int) repository.WriterFunc
}

func New(repo Repo, pool *Pool) *CommandService {
	return &CommandService{
		Repo:    repo,
		Pool:    pool,
		Workers: sync.Map{},
	}
}

func (cs *CommandService) CreateCommand(cmd *model.Command, isLong bool) error {
	slog.Info("in the service.CreateCommand")
	ctx := context.Background()

	err := cs.Repo.Create(ctx, cmd)
	if err != nil {
		return err
	}

	cs.Pool.Go(func() {
		slog.Info("in the func of the Pool.Go")
		err := cs.ExecCommand(ctx, *cmd, isLong)
		if err != nil {
			slog.ErrorContext(ctx, "failed to exec the command", err)
		}
	})

	return nil
}

func (cs *CommandService) ExecCommand(ctx context.Context, cmd model.Command, isLong bool) (err error) {
	var wg sync.WaitGroup
	errch := make(chan error, 1)

	defer func() {
		wg.Wait()
		close(errch)
		receivedErr := <-errch
		if receivedErr != nil {
			err = errors.Join(err, cs.Repo.SetError(ctx, cmd.Id, receivedErr.Error()))
			slog.Error("failed to execute command", "error", receivedErr.Error())
		} else {
			if err = cs.Repo.SetSuccess(ctx, cmd.Id); err != nil {
				slog.Error("failed to set command success in db", "error", err.Error())
			}
		}
	}()

	slog.Info("in the func ExecCommand")
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

	_, err = cs.Repo.Writer(ctx, cmd.Id)([]byte(accumulate))
	if err != nil {
		slog.Error("failed to write cmd result into db", "error", err.Error())
		return
	}
}

func (cs *CommandService) long(ctx context.Context, stdout io.ReadCloser, scanner *bufio.Scanner, execCmd *exec.Cmd, errch chan error, wg *sync.WaitGroup, cmd *model.Command) {
	defer stdout.Close()
	defer wg.Done()

	var err error
	defer func() {
		errch <- err
	}()

	for scanner.Scan() {
		line := scanner.Text()
		_, err := cs.Repo.Writer(ctx, cmd.Id)([]byte(line + "\n"))
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
