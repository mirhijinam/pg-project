package service

import (
	"context"
	"errors"
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
		err := cs.ExecCommand(ctx, *cmd)
		if err != nil {
			slog.ErrorContext(ctx, "failed to exec the command", err)
		}
	})

	return nil
}

func (cs *CommandService) ExecCommand(ctx context.Context, cmd model.Command) (err error) {

	defer func() {
		if err == nil {
			err = cs.Repo.SetSuccess(ctx, cmd.Id)
			slog.Info("successfully executed and wrote into db", slog.String("cmd", cmd.Raw))
		} else {
			err = errors.Join(err, cs.Repo.SetError(ctx, cmd.Id, err.Error()))
		}
	}()

	slog.Info("in the func ExecCommand")
	execCmd := exec.Command("/bin/sh", "-c", cmd.Raw)

	execCmd.Stdout = cs.Repo.Writer(ctx, cmd.Id)

	cs.Workers.Store(cmd.Id, &model.CommandExec{
		Cmd:  cmd,
		Exec: execCmd,
	})

	return execCmd.Run()
}
