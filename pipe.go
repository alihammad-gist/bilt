package main

import (
	"io"
	"os/exec"
)

func Pipe(r io.Reader, w io.Writer, cmds ...*exec.Cmd) error {

	var lcmd *exec.Cmd
	for _, cmd := range cmds {
		if lcmd != nil {
			if stdout, err := lcmd.StdoutPipe(); err == nil {
				cmd.Stdin = stdout
			} else {
				errlogger.Fatal(err)
			}
			if err := lcmd.Start(); err != nil {
				errlogger.Fatal(err)
			}
		} else {
			cmd.Stdin = r
		}
		lcmd = cmd
	}

	lcmd.Stdout = w
	if err := lcmd.Start(); err != nil {
		errlogger.Fatal(err)
	}

	for _, cmd := range cmds {
		if err := cmd.Wait(); err != nil {
			errlogger.Fatal(err)
		}
	}

	return nil
}
