package main

import (
	"errors"
	"log"

	"github.com/alihammad-gist/sniffy"
)

var (
	ErrNoDirProvided = errors.New("Suite doesn't have any directories")
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

/*
[
	{
		"dirs": ["/path/to/dir", "/home/ali"],
		"exts": [".php", ".css"],
		"cmds": ["browserify"],
		"src" : "main.js",
		"dest": "dest.js"
	},
	{
		"dirs": ["/", "/usr"],
		"exts": [".bat", ".bin"],
		"cmds": ["cat"],
		"src" : "main.sh",
		"dest": "dest.sh"
	}
]
*/

func main() {
	suites, err := getSuites()
	if err != nil {
		log.Fatal(err)
	}

	// getting EventTransmitters
	// watching for events
	var trans []*sniffy.EventTransmitter
	for _, s := range suites {
		t, err := getTransmitter(s)
		if err != nil {
			log.Fatal(err)
		}
		go func() {
			for e := range t.Events {
				// TODO:
				// Set Suite.Root to os/exec.Cmd.Dir (if not empty)
				//
			}
		}()
		trans = append(trans, t)
	}
}

func getTransmitter(s *Suite) (*sniffy.EventTransmitter, error) {
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

func getSuites() []*Suite {
	// Turn paths to absolute paths include Suite.Src and Suite.Dest
}
