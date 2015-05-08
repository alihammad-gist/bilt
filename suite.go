package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/alihammad-gist/sniffy"
)

var (
	ErrNoDirProvided = errors.New("Suite doesn't have any directories")
	ErrNoSrcOrDest   = errors.New("Dest or source not provided for one of the suites")
)

// Looks if config file is passed as argument in the command-line
// or if it is present in current working directory, If none it will return
// an error. First returned value is nil if error isn't
//
// If there is error Unmarshalling json, error will be returned.
func Suites() ([]*Suite, error) {
	var (
		suites []*Suite
		d      []byte
		err    error
	)

	if len(os.Args) == 2 {
		d, err = ioutil.ReadFile(os.Args[1])
		if err != nil {
			return nil, err
		}
	} else {
		d, err = ioutil.ReadFile(DEFAULT_CONFIG)
		if err != nil {
			return nil, err
		}
	}

	if err = json.Unmarshal(d, &suites); err != nil {
		return nil, err
	}

	for _, s := range suites {
		if len(s.Dirs) == 0 {
			return nil, ErrNoDirProvided
		}

		if s.Src == "" || s.Dest == "" {
			return nil, ErrNoSrcOrDest
		}

		if err = AbsPaths(s); err != nil {
			return nil, err
		}
	}

	return suites, nil
}

// Absolutizes paths inside a suite
func AbsPaths(s *Suite) error {

	// Watch paths
	var (
		p   string
		err error
	)
	for i, _ := range s.Dirs {
		p = s.Dirs[i]
		if !filepath.IsAbs(p) {
			if s.Root == "" {
				s.Dirs[i], err = filepath.Abs(p)
				if err != nil {
					return err
				}
			} else {
				s.Dirs[i] = filepath.Join(s.Root, p)
			}
		}
	}

	// Source and destination paths
	if s.Root == "" {
		s.Src, err = filepath.Abs(s.Src)
		if err != nil {
			return err
		}
		pwd, err := os.Getwd()
		if err != nil {
			return err
		}
		s.Dest = filepath.Join(pwd, s.Dest)

		// setting root for CMD
		if rt, err := os.Getwd(); err != nil {
			errlogger.Fatal(err)
		} else {
			s.Root = rt
		}
	} else {
		s.Src = filepath.Join(s.Root, s.Src)
		s.Dest = filepath.Join(s.Root, s.Dest)
	}

	return nil
}

func (s *Suite) Exec() error {
	var cmds []*exec.Cmd
	for _, cstr := range s.Cmds {
		cparts := strings.Split(cstr, " ")
		c := exec.Command(cparts[0], cparts[1:]...)
		c.Dir = s.Root
		cmds = append(cmds, c)
	}

	// check if src file is present
	if _, err := os.Stat(s.Src); os.IsNotExist(err) {
		errlogger.Fatalf("no such file or directory: %s", s.Src)
	}

	r, err := os.Open(s.Src)
	if err != nil {
		return err
	}
	defer r.Close()

	w, err := os.Create(s.Dest)
	if err != nil {
		return err
	}
	defer w.Close()

	return Pipe(r, w, os.Stderr, cmds...)
}

func (s *Suite) Transmitter() (*sniffy.EventTransmitter, error) {
	var filters []sniffy.Filter

	filters = append(
		filters,
		sniffy.ChildFilter(s.Dirs...),
		sniffy.ExcludePathFilter(s.Dest),
		sniffy.TooSoonFilter(time.Second),
	)

	if len(s.Exts) > 0 {
		filters = append(filters, sniffy.ExtFilter(s.Exts...))
	}
	return sniffy.Transmitter(filters...), nil
}
