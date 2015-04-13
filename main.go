package main

import (
	"encoding/json"
	"fmt"
	"log"
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

func main() {
	jsonBlob := []byte(`
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
	`)
	var conf []Suite
	err := json.Unmarshal(jsonBlob, &conf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v", conf)
}
