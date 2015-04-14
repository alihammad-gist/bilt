package main

import (
	"errors"
	"io"
	"os/exec"
)

var (
	ErrNoCommandsProvided = errors.New("No commands were provided")
)

func Pipe(in io.Reader, out io.Writer, cmds ...*exec.Cmd) error {
	if len(cmds) == 0 {
		return ErrNoCommandsProvided
	}

	// setting stdin and out
	cmds[0].Stdin = in
	cmds[len(cmds)-1].Stdout = out

	for i, _ := range cmds {
		if cmds[i].Stdin == nil {
			pout, err := cmds[i-1].StdoutPipe()
			if err != nil {
				return err
			}
			cmds[i].Stdin = pout
		}
	}

	// starting
	for i, _ := range cmds {
		if err := cmds[i].Start(); err != nil {
			return err
		}
	}

	// waiting
	for i, _ := range cmds {
		if err := cmds[i].Wait(); err != nil {
			return err
		}
	}

	return nil

}
