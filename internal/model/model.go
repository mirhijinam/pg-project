package model

import (
	"database/sql"
	"os/exec"
	"time"
)

type Command struct {
	Id        int
	Name      string
	Raw       string
	Status    sql.NullString
	ErrorMsg  sql.NullString
	Logs      sql.NullString
	CreatedAt time.Time
	UpdatedAt *time.Time
}

type CommandExec struct {
	Cmd  Command
	Exec *exec.Cmd
}
