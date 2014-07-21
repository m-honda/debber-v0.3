package main

import (
	"log"
	"os"
)
const (
	TaskInit = "init"
	TaskGenChangelog = "gen:changelog"
	TaskGenDeb = "gen:deb"
	TaskGenDev = "gen:dev"
	TaskGenSource = "gen:source"
	TaskDebControl = "deb:control"
	TaskDebContents = "deb:contents"
	TaskDebContentsDebian = "deb:contents-debian"
)

var tasks = []string{
	TaskInit,
	TaskGenChangelog,
	TaskGenDeb,
	TaskGenDev,
	TaskGenSource,
	TaskDebControl,
	TaskDebContents,
	TaskDebContentsDebian,
}

func main() {
	name := "debber"
	log.SetPrefix("[" + name + "] ")
	if len(os.Args) < 2 {
		log.Fatalf("Please specify a task (one of: %v)", tasks)
	}
	task := os.Args[1]
	args := os.Args[2:]
	//var debFile string
	//fs.StringVar(&debFile, "file", "", ".deb file")
	switch task {
	case TaskInit:
		initDebber(args)
	case TaskGenChangelog:
		genChangelog(args)
	case TaskGenDeb:
		genDeb(args)
	case TaskGenDev:
		genDev(args)
	case TaskGenSource:
		genSource(args)
	case TaskDebControl:
		debControl(args)
	case TaskDebContents:
		debContents(args)
	case TaskDebContentsDebian:
		debContentsDebian(args)
	default:
		log.Printf("Unrecognised command '%s'", task)
		log.Fatalf("Please specify one of: %v", tasks)
	}

}
