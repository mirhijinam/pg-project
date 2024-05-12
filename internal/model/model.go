package model

import (
	"errors"
	"os/exec"
	"time"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Command struct {
	Id        int        `json:"id"`
	Name      string     `json:"name"`
	Raw       string     `json:"raw"`
	Status    string     `json:"status,omitempty"`
	ErrorMsg  string     `json:"error_msg,omitempty"`
	Logs      string     `json:"logs,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

type CommandExec struct {
	Cmd  Command
	Exec *exec.Cmd
}
