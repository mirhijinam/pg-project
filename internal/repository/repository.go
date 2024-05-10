package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/mirhijinam/pg-project/internal/model"
)

type CommandRepo struct {
	db *sql.DB
}

func New(db *sql.DB) *CommandRepo {
	return &CommandRepo{
		db: db,
	}
}

const (
	cmdStatusSuccess = "success"
	cmdStatusStopped = "stopped"
	cmdStatusError   = "error"
)

func (cr *CommandRepo) Create(ctx context.Context, cmd *model.Command) (err error) {
	stmt := `
		INSERT INTO commands (name, raw)
		VALUES ($1, $2)
		RETURNING id, created_at`
	args := []interface{}{cmd.Name, cmd.Raw}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	return cr.db.QueryRowContext(ctx, stmt, args...).Scan(&cmd.Id, &cmd.CreatedAt)
}

func (cr *CommandRepo) SetError(ctx context.Context, id int, errMsg string) (err error) {
	stmt := `
			UPDATE commands
			SET status = $1, error_msg = $2, updated_at = now()
			WHERE id = $3`

	if errMsg == "signal: killed" {
		_, err = cr.db.ExecContext(ctx, stmt, cmdStatusStopped, errMsg, id)
	} else {
		_, err = cr.db.ExecContext(ctx, stmt, cmdStatusError, errMsg, id)
	}

	return err
}

func (cr *CommandRepo) SetSuccess(ctx context.Context, id int) (err error) {
	stmt := `
		UPDATE commands 
		SET status = $1, updated_at = now() 
		WHERE id = $2`

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	result, err := cr.db.ExecContext(ctx, stmt, cmdStatusSuccess, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("record not found")
	}

	return nil
}

func (cr *CommandRepo) GetList(ctx context.Context) (cmdList []model.Command, err error) {
	stmt := `
		SELECT c.id, c.name, c.raw, c.status, c.error_msg, cl.logs, c.created_at, c.updated_at 
		FROM commands c
		LEFT JOIN command_logs cl ON c.id = cl.command_id`
	rowList, err := cr.db.QueryContext(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rowList.Close()

	for rowList.Next() {
		var cmd model.Command
		err := rowList.Scan(
			&cmd.Id,
			&cmd.Name,
			&cmd.Raw,
			&cmd.Status,
			&cmd.ErrorMsg,
			&cmd.Logs,
			&cmd.CreatedAt,
			&cmd.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		cmdList = append(cmdList, cmd)
	}
	return cmdList, nil
}

func (cr *CommandRepo) GetOne(ctx context.Context, id int) (cmd model.Command, err error) {
	stmt := `
        SELECT c.id, c.name, c.raw, c.status, c.error_msg, cl.logs, c.created_at, c.updated_at
        FROM commands c
        LEFT JOIN command_logs cl ON c.id = cl.command_id
        WHERE c.id = $1`

	row := cr.db.QueryRowContext(ctx, stmt, id)
	err = row.Scan(
		&cmd.Id,
		&cmd.Name,
		&cmd.Raw,
		&cmd.Status,
		&cmd.ErrorMsg,
		&cmd.Logs,
		&cmd.CreatedAt,
		&cmd.UpdatedAt,
	)
	if err != nil {
		return model.Command{}, err
	}
	return cmd, nil
}

func (cr *CommandRepo) Writer(ctx context.Context, id int) WriterFunc {
	stmt := `
		INSERT INTO command_logs (command_id, logs)
		VALUES ($1, $2)
		ON CONFLICT (command_id) DO UPDATE 
		SET logs = command_logs.logs || $2`

	return func(p []byte) (n int, err error) {
		_, err = cr.db.ExecContext(ctx, stmt, id, p)
		if err != nil {
			return 0, err
		}

		return len(p), nil
	}
}

type WriterFunc func(p []byte) (n int, err error)

func (w WriterFunc) Write(p []byte) (n int, err error) {
	return w(p)
}
