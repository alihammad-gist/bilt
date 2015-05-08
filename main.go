package main

import (
	"log"
	"os"

	"github.com/alihammad-gist/sniffy"
)

const (
	DEFAULT_CONFIG = "bilt.json"
)

type (
	Suite struct {
		Dirs  []string `json:"dirs"`
		Exts  []string `json:"exts"`
		Root  string   `json:"root"`
		Cmds  []string `json:"cmds"`
		Src   string   `json:"src"`
		Dest  string   `json:"dest"`
		Label string   `json:"label"`
	}
)

var (
	errlogger *log.Logger
	evlogger  *log.Logger
)

func init() {
	errlogger = log.New(os.Stderr, "!  ", log.Lshortfile)
	evlogger = log.New(os.Stdout, "-> ", log.Ltime)
}

func main() {
	suites, err := Suites()
	if err != nil {
		errlogger.Fatal(err)
	}

	// getting EventTransmitters
	// watching for events
	var trans []*sniffy.EventTransmitter
	for _, s := range suites {
		t, err := s.Transmitter()
		if err != nil {
			errlogger.Fatal(err)
		}
		go func() {
			for e := range t.Events {
				evlogger.Println(e.Name)
				if err := s.Exec(); err != nil {
					errlogger.Println(err)
				}
				evlogger.Println("Exec Done (", s.Label, ")")
			}
		}()
		trans = append(trans, t)
	}

	w, err := sniffy.NewWatcher(trans...)
	if err != nil {
		errlogger.Println(err)
	}

	// watching for errors
	go func() {
		for err := range w.Errors {
			errlogger.Println(err)
		}
	}()

	for _, s := range suites {
		for _, d := range s.Dirs {
			err = w.AddDir(d)
			if err != nil {
				errlogger.Println(err, d)
			}
		}
	}

	done := make(chan bool)
	<-done
}
