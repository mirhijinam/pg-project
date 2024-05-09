package model

import (
	"os/exec"
	"time"
)

type Command struct {
	Id        int
	Name      string
	Raw       string
	CreatedAt time.Time
	UpdatedAt *time.Time
}

type CommandExec struct {
	Cmd  Command
	Exec *exec.Cmd
}
