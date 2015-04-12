package main

import (
	"log"
	"os"
	"os/exec"
)

type (
	Pipe struct {
		cmds []*exec.Cmd
	}

	ConfigDirs struct {
		Path
	}

	Config struct {
	}
)

func main() {
	// cat := exec.Command("cat", "file.txt")
	// stdout, _ := cat.StdoutPipe()

	// grep := exec.Command("grep", "flld")
	// grep.Stdin = stdout
	// grep.Stdout = os.Stdout

	// cat.Start()
	// grep.Start()
	// cat.Wait()
	// grep.Wait()

	p, err := NewPipe(
		exec.Command("cat", "file.txt"),
		exec.Command("grep", "flld"),
		exec.Command("grep", "sfds"),
	)

	if err != nil {
		log.Fatal(err)
	}

	if err := p.Exec(); err != nil {
		log.Fatal(err)
	}
	if err := p.Exec(); err != nil {
		log.Fatal("Rerun1: ", err)
	}
	if err := p.Exec(); err != nil {
		log.Fatal(err)
	}
	if err := p.Exec(); err != nil {
		log.Fatal(err)
	}
}

func NewPipe(cmds ...*exec.Cmd) (*Pipe, error) {

	p := &Pipe{
		cmds: cmds,
	}

	var lcmd *exec.Cmd
	for _, cmd := range p.cmds {
		if lcmd != nil {
			if stdout, err := lcmd.StdoutPipe(); err == nil {
				cmd.Stdin = stdout
			} else {
				return nil, err
			}
		}
		lcmd = cmd
	}
	lcmd.Stdout = os.Stdout

	return p, nil
}

func (p *Pipe) Exec() error {

	for _, cmd := range p.cmds {
		if err := cmd.Start(); err != nil {
			return err
		}
	}

	for _, cmd := range p.cmds {
		if err := cmd.Wait(); err != nil {
			return err
		}
	}

	return nil
}
