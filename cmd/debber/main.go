package main

import (
	"log"
	"os"
)
const (
	TaskInit = "init"
	TaskChangelogAdd = "changelog:add"
	TaskGenDeb = "deb:gen"
	TaskDebControl = "deb:control"
	TaskDebContents = "deb:contents"
	TaskDebContentsDebian = "deb:contents-debian"
	TaskGenDev = "dev:gen"
	TaskGenSource = "source:gen"
)

var tasks = []string{
	TaskInit,
	TaskChangelogAdd,
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
	switch task {
	case TaskInit:
		initDebber(args)
	case TaskChangelogAdd:
		changelogAdd(args)
	case TaskGenDeb:
		debGen(args)
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
