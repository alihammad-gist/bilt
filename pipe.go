package main

import (
	"errors"
	"io"
	"os/exec"
)

var (
	ErrTooFewCmds = errors.New("Function pipe expects at least two commands")
)

func Pipe(w io.Writer, cmds ...*exec.Cmd) error {
	if len(cmds) < 2 {
		return ErrTooFewCmds
	}

	var lcmd *exec.Cmd
	for _, cmd := range cmds {
		if lcmd != nil {
			if stdout, err := lcmd.StdoutPipe(); err == nil {
				cmd.Stdin = stdout
			} else {
				return err
			}
			if err := lcmd.Start(); err != nil {
				return err
			}
		}
		lcmd = cmd
	}

	lcmd.Stdout = w
	if err := lcmd.Start(); err != nil {
		return err
	}

	for _, cmd := range cmds {
		if err := cmd.Wait(); err != nil {
			return err
		}
	}

	return nil
}
