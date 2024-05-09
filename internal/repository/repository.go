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

func (cr *CommandRepo) SetSuccess(ctx context.Context, id int) (err error) {
	stmt := `
	UPDATE commands 
	SET status = 'success', updated_at = now() 
	WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	result, err := cr.db.ExecContext(ctx, stmt, id)
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

func (cr *CommandRepo) SetError(ctx context.Context, id int, errMsg string) (err error) {
	stmt := `
	UPDATE commands
	SET status = 'error', error_msg = $1, updated_at = now()
	WHERE id = $2`

	_, err = cr.db.ExecContext(ctx, stmt, errMsg, id)

	return err
}

func (cr *CommandRepo) List(ctx context.Context) (cmds []model.Command, err error) {

	return nil, nil
}
func (cr *CommandRepo) Get(ctx context.Context, id string) (cmd model.Command, err error) {

	return model.Command{}, nil
}

func (cr *CommandRepo) Writer(ctx context.Context, id int) WriterFunc {
	stmt := `
		INSERT INTO command_logs(command_id, logs)
		VALUES ($1, $2)`

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
