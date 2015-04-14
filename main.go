package main

import (
	"log"

	"github.com/alihammad-gist/sniffy"
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

func main() {
	suites, err := Suites()
	if err != nil {
		log.Fatal(err)
	}

	// getting EventTransmitters
	// watching for events
	var trans []*sniffy.EventTransmitter
	for _, s := range suites {
		t, err := s.Transmitter()
		if err != nil {
			log.Fatal(err)
		}
		go func() {
			for e := range t.Events {
				log.Println(e)
				if err := s.Exec(); err != nil {
					log.Println(err)
				}
			}
		}()
		trans = append(trans, t)
	}

	w, err := sniffy.NewWatcher(trans...)
	if err != nil {
		log.Fatal(err)
	}

	// watching for errors
	go func() {
		for err := range w.Errors {
			log.Println(err)
		}
	}()

	for _, s := range suites {
		for _, d := range s.Dirs {
			err = w.AddDir(d)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	done := make(chan bool)
	<-done
}
