package main

import (
	"log"
	"os"
)

var tasks = []string {
	"init",
	"gen:deb",
	"deb:contents",
	"deb:control",
}

func main() {
	name := "debber"
	log.SetPrefix("[" + name + "] ")
	if len(os.Args)<2 {
		log.Fatalf("Please specify a task (one of: %v)", tasks)
	}
	task := os.Args[1]
	args := os.Args[2:]
	//var debFile string
	//fs.StringVar(&debFile, "file", "", ".deb file")
	switch task {
	case "init":
		initDebber(args)
	case "gen:changelog": 
		genChangelog(args)
	case "deb:control": 
		debControl(args)

	case "deb:contents":
		debContents(args)
	case "deb:contents-debian":
		debContentsDebian(args)
	default:
		log.Fatalf("No command specified")
	}

}
