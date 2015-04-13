package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alihammad-gist/sniffy"
)

// Looks if config file is passed as argument in the command-line
// or if it is present in current working directory
func Suites() []*Suite {
	var suites []*Suite
	// Turn paths to absolute paths include Suite.Src and Suite.Dest

	for _, s := range suites {

	}
}

func AbsPaths(s *Suite) error {
	// Watch paths
	var (
		p   string
		err error
	)
	for i, _ := range dirs {
		p = dirs[i]
		if !filepath.IsAbs(p) {
			if s.Root == "" {
				dirs[i], err = filepath.Abs(p)
				if err != nil {
					return err
				}
			} else {
				dirs[i] = filepath.Join(s.Root, p)
			}
		}
	}

	// Source and destination paths
	if s.Root == "" {
		s.Src, err = filepath.Abs(s.Src)
		if err != nil {
			return err
		}
		s.Dest, err = filepath.Abs(s.Dest)
		if err != nil {
			return err
		}
	} else {
		s.Src = filepath.Join(s.Root, s.Src)
		s.Dest = filepath.Join(s.Root, s.Dest)
	}
}

func (s *Suite) Exec() error {
	var cmds []*exec.Cmd
	for _, cstr := range s.Cmds {
		cparts := strings.Split(cstr)
		c := exec.Command(cparts[0], cparts[1:]...)
		c.Dir = s.Root
		cmds = append(cmds, c)
	}
	r, err := os.Open(s.Src)
	if err != nil {
		return err
	}
	w, err := os.OpenFile(s.Dest, os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	return Pipe(r, w, cmds...)
}

func (s *Suite) Transmitter() (*sniffy.EventTransmitter, error) {
	if len(s.Dirs) == 0 {
		return nil, ErrNoDirProvided
	}
	var filters []sniffy.Filter

	filters = append(filters, sniffy.PathFilter(s.Dirs...))

	if len(s.Exts) > 0 {
		filters = append(filters, sniffy.ExtFilter(s.Exts...))
	}
	return &sniffy.Transmitter(filters...), nil
}
