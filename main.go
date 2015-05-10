package main

import (
	"log"
	"os"

	"github.com/jaschaephraim/lrserver"

	"github.com/alihammad-gist/sniffy"
)

const (
	DEFAULT_CONFIG = "bilt.json"

	LIVERLD_EVENT = "?updated"
	LIVERLD_DEST  = "?dest"
)

type (
	Suite struct {
		Dirs    []string `json:"dirs"`
		Exts    []string `json:"exts"`
		Root    string   `json:"root"`
		Cmds    []string `json:"cmds"`
		Src     string   `json:"src"`
		Dest    string   `json:"dest"`
		Label   string   `json:"label"`
		LiveRld string   `json:"liveReload"`

		trans *sniffy.EventTransmitter
	}
)

var (
	errlogger    *log.Logger
	evlogger     *log.Logger
	liveRld      *lrserver.Server
	runSuiteChan chan *Suite
)

func init() {
	var err error
	errlogger = log.New(os.Stderr, "[Error]  ", log.Lshortfile)
	evlogger = log.New(os.Stdout, "[BILT] ", log.Ltime)

	// live reload server
	liveRld, err = lrserver.New(lrserver.DefaultName, lrserver.DefaultPort)
	if err != nil {
		errlogger.Fatal(err)
	}

	go func() {
		if err = liveRld.ListenAndServe(); err != nil {
			errlogger.Fatal(err)
		}
	}()
}

func main() {
	runSuiteChan = make(chan *Suite)
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
				runSuiteChan <- s
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

	// watching for suites
	go func() {
		for s := range runSuiteChan {
			if err := s.Exec(); err != nil {
				errlogger.Println(err)
			} else {
				s.Publish(liveRld)
			}
			evlogger.Println("Exec Done (", s.Label, ")")
		}
	}()

	// adding dirs
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
