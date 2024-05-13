package repository

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	mockSql "github.com/DATA-DOG/go-sqlmock"
	"github.com/mirhijinam/pg-project/internal/model"
)

func TestCommandRepo_Create(t *testing.T) {
	db, mock, err := mockSql.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	type fields struct {
		db *sql.DB
	}
	type args struct {
		ctx context.Context
		cmd *model.Command
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		setup   func()
	}{
		{
			name: "Successful create",
			fields: fields{
				db: db,
			},
			args: args{
				ctx: context.Background(),
				cmd: &model.Command{Name: "echo", Raw: "echo 'test'"},
			},
			wantErr: false,
			setup: func() {
				rows := mockSql.NewRows([]string{"id", "created_at"}).AddRow(1, time.Now())
				mock.ExpectQuery(`^INSERT INTO commands \(name, raw\) VALUES \(\$1, \$2\) RETURNING id, created_at$`).WithArgs("echo", "echo 'test'").WillReturnRows(rows)
			},
		},
		{
			name: "Failure success update",
			fields: fields{
				db: db,
			},
			args: args{
				ctx: context.Background(),
				cmd: &model.Command{Name: "echo", Raw: "echo 'test'"},
			},
			wantErr: true,
			setup: func() {
				mock.ExpectQuery(`^INSERT INTO commands \(name, raw\) VALUES \(\$1, \$2\) RETURNING id, created_at$`).WithArgs("echo", "echo 'test'").WillReturnError(errors.New("insert error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			cr := &CommandRepo{
				db: tt.fields.db,
			}
			if err := cr.Create(tt.args.ctx, tt.args.cmd); (err != nil) != tt.wantErr {
				t.Errorf("CommandRepo.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCommandRepo_SetError(t *testing.T) {
	db, mock, err := mockSql.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	type fields struct {
		db *sql.DB
	}
	type args struct {
		ctx    context.Context
		id     int
		errMsg string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		setup   func()
	}{
		{
			name: "Successful error update",
			fields: fields{
				db: db,
			},
			args: args{
				ctx:    context.TODO(),
				id:     1,
				errMsg: "error occurred",
			},
			wantErr: false,
			setup: func() {
				mock.ExpectExec(`UPDATE commands SET status = \$1, error_msg = \$2, updated_at = now\(\) WHERE id = \$3`).
					WithArgs(cmdStatusError, "error occurred", 1).
					WillReturnResult(mockSql.NewResult(1, 1))
			},
		},
		{
			name: "Failure error update",
			fields: fields{
				db: db,
			},
			args: args{
				ctx:    context.TODO(),
				id:     1,
				errMsg: "error occurred",
			},
			wantErr: true,
			setup: func() {
				mock.ExpectExec(`UPDATE commands SET status = \$1, error_msg = \$2, updated_at = now\(\) WHERE id = \$3`).
					WithArgs(cmdStatusError, "error occurred", 1).
					WillReturnError(errors.New("database error"))
			},
		},
		{
			name: "Case with stopped command",
			fields: fields{
				db: db,
			},
			args: args{
				ctx:    context.TODO(),
				id:     2,
				errMsg: "signal: killed",
			},
			wantErr: false,
			setup: func() {
				mock.ExpectExec(`UPDATE commands SET status = \$1, error_msg = \$2, updated_at = now\(\) WHERE id = \$3`).
					WithArgs(cmdStatusStopped, "signal: killed", 2).
					WillReturnResult(mockSql.NewResult(1, 1))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			cr := &CommandRepo{
				db: tt.fields.db,
			}
			if err := cr.SetError(tt.args.ctx, tt.args.id, tt.args.errMsg); (err != nil) != tt.wantErr {
				t.Errorf("CommandRepo.SetError() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCommandRepo_SetSuccess(t *testing.T) {
	db, mock, err := mockSql.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	type fields struct {
		db *sql.DB
	}
	type args struct {
		ctx context.Context
		id  int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		setup   func()
	}{
		{
			name: "Successful update",
			fields: fields{
				db: db,
			},
			args: args{
				ctx: context.TODO(),
				id:  1,
			},
			wantErr: false,
			setup: func() {
				mock.ExpectExec(`UPDATE commands SET status = \$1, updated_at = now\(\) WHERE id = \$2`).
					WithArgs(cmdStatusSuccess, 1).
					WillReturnResult(mockSql.NewResult(0, 1)) // 1 row affected
			},
		},
		{
			name: "No rows affected",
			fields: fields{
				db: db,
			},
			args: args{
				ctx: context.TODO(),
				id:  1,
			},
			wantErr: true,
			setup: func() {
				mock.ExpectExec(`UPDATE commands SET status = \$1, updated_at = now\(\) WHERE id = \$2`).
					WithArgs(cmdStatusSuccess, 1).
					WillReturnResult(mockSql.NewResult(0, 0)) // 0 rows affected
			},
		},
		{
			name: "Database error",
			fields: fields{
				db: db,
			},
			args: args{
				ctx: context.TODO(),
				id:  1,
			},
			wantErr: true,
			setup: func() {
				mock.ExpectExec(`UPDATE commands SET status = \$1, updated_at = now\(\) WHERE id = \$2`).
					WithArgs(cmdStatusSuccess, 1).
					WillReturnError(errors.New("database error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			cr := &CommandRepo{
				db: tt.fields.db,
			}
			if err := cr.SetSuccess(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("CommandRepo.SetSuccess() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCommandRepo_GetList(t *testing.T) {
	db, mock, err := mockSql.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	type fields struct {
		db *sql.DB
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []model.Command
		wantErr bool
	}{
		{
			name: "Successful fetch",
			fields: fields{
				db: db,
			},
			args: args{
				ctx: context.TODO(),
			},
			want: []model.Command{
				{Id: 1, Name: "Command One", Raw: "echo 'hello'", Status: "", ErrorMsg: "", Logs: ""},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cr := &CommandRepo{
				db: tt.fields.db,
			}

			rows := mockSql.NewRows([]string{"id", "name", "raw", "status", "error_msg", "logs", "created_at", "updated_at"}).
				AddRow(1, "Command One", "echo 'hello'", "", "", "", time.Now(), time.Now())
			mock.ExpectQuery(`SELECT`).WillReturnRows(rows)

			got, err := cr.GetList(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("CommandRepo.GetList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for i, gotCmd := range got {
				wantCmd := tt.want[i]
				if gotCmd.Id != wantCmd.Id || gotCmd.Name != wantCmd.Name || gotCmd.Raw != wantCmd.Raw ||
					gotCmd.Status != wantCmd.Status || gotCmd.ErrorMsg != wantCmd.ErrorMsg || gotCmd.Logs != wantCmd.Logs {
					t.Errorf("CommandRepo.GetList() = %v, want %v at index %d", gotCmd, wantCmd, i)
				}
			}
		})
	}
}

func TestCommandRepo_GetOne(t *testing.T) {
	db, mock, err := mockSql.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	type fields struct {
		db *sql.DB
	}
	type args struct {
		ctx context.Context
		id  int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    model.Command
		wantErr bool
	}{
		{
			name: "Successful fetch",
			fields: fields{
				db: db,
			},
			args: args{
				ctx: context.TODO(),
				id:  1,
			},
			want:    model.Command{Id: 1, Name: "Command One", Raw: "echo 'hello'"},
			wantErr: false,
		},
		{
			name: "Not found",
			fields: fields{
				db: db,
			},
			args: args{
				ctx: context.TODO(),
				id:  1,
			},
			want:    model.Command{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cr := &CommandRepo{
				db: tt.fields.db,
			}

			if tt.name == "Successful fetch" {
				rows := sqlmock.NewRows([]string{"id", "name", "raw", "status", "error_msg", "logs", "created_at", "updated_at"}).
					AddRow(1, "Command One", "echo 'hello'", "", "", "", time.Now(), time.Now())
				mock.ExpectQuery(`SELECT`).WithArgs(1).WillReturnRows(rows)
			} else if tt.name == "Not found" {
				mock.ExpectQuery(`SELECT`).WithArgs(1).WillReturnError(sql.ErrNoRows)
			}

			got, err := cr.GetOne(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("CommandRepo.GetOne() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Compare non-time fields manually
			if got.Id != tt.want.Id || got.Name != tt.want.Name || got.Raw != tt.want.Raw ||
				got.Status != tt.want.Status || got.ErrorMsg != tt.want.ErrorMsg || got.Logs != tt.want.Logs {
				t.Errorf("CommandRepo.GetOne() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommandRepo_Writer(t *testing.T) {
	db, mock, err := mockSql.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	type fields struct {
		db *sql.DB
	}
	type args struct {
		ctx context.Context
		id  int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		data    []byte
	}{
		{
			name: "Successful write",
			fields: fields{
				db: db,
			},
			args: args{
				ctx: context.TODO(),
				id:  1,
			},
			data:    []byte("log data"),
			wantErr: false,
		},
		{
			name: "Database error",
			fields: fields{
				db: db,
			},
			args: args{
				ctx: context.TODO(),
				id:  1,
			},
			data:    []byte("log data"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cr := &CommandRepo{
				db: tt.fields.db,
			}

			writer := cr.Writer(tt.args.ctx, tt.args.id)
			if tt.name == "Successful write" {
				mock.ExpectExec(`INSERT INTO command_logs`).WithArgs(tt.args.id, tt.data).WillReturnResult(mockSql.NewResult(1, 1))
			} else if tt.name == "Database error" {
				mock.ExpectExec(`INSERT INTO command_logs`).WithArgs(tt.args.id, tt.data).WillReturnError(errors.New("db error"))
			}

			_, err := writer(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("CommandRepo.Writer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
