package main

import (
	"errors"

	"github.com/alihammad-gist/sniffy"
)

var (
	ErrNoDirProvided = errors.New("Suite don't have any directory")
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
		"exts": ["php", "css"],
		"cmds": ["browserify"],
		"src" : "main.js",
		"dest": "dest.js"
	},
	{
		"dirs": ["/", "/usr"],
		"exts": ["bat", "bin"],
		"cmds": ["cat"],
		"src" : "main.sh",
		"dest": "dest.sh"
	}
]
*/

func main() {

}

func getTransmitter(s *Suite) (*sniffy.EventTransmitter, error) {

}
