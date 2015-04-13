package main

import (
	"log"
	"os"
	"os/exec"
)

type (
	Suite struct {
		Dirs []string `json:"dirs"`
		Exts []string `json:"exts"`
		Cmds []string `json:"cmds"`
		Src  string   `json:"src"`
		Dest string   `json:"dest"`
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
	f, err := os.Open("file.txt")
	defer f.Close()

	if err != nil {
		log.Fatal(err)
	}
	err = Pipe(
		f,
		os.Stdout,
		exec.Command("grep", "asd"),
		exec.Command("grep", ";"),
	)
	if err != nil {
		log.Fatal(err)
	}
}
