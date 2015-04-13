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
		t, err := s.Transmitter()
		if err != nil {
			log.Fatal(err)
		}
		go func() {
			for e := range t.Events {
				if err := s.Exec(); err != nil {
					log.Println(err)
				}
			}
		}()
		trans = append(trans, t)
	}
}
