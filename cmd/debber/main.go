package main

import (
	"log"
	"os"
)

const (
	TaskInit              = "init"
	TaskChangelogAdd      = "changelog:add"
	TaskGenDeb            = "deb:gen"
	TaskDebControl        = "deb:control"
	TaskDebContents       = "deb:contents"
	TaskDebContentsDebian = "deb:contents-debian"
	TaskGenDev            = "dev:gen"
	TaskGenSource         = "source:gen"
)

var tasks = []string{
	TaskInit,
	TaskChangelogAdd,
	TaskGenDeb,
	TaskGenSource,
	TaskDebControl,
	TaskDebContents,
	TaskDebContentsDebian,
}

func main() {
	name := "debber"
	log.SetPrefix("[" + name + "] ")
	if len(os.Args) < 2 {
		log.Printf("Please specify a task (one of: %v)", tasks)
		log.Printf("For help on any task, use `debber <task> -h`")
		os.Exit(1)
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
	case TaskGenSource:
		sourceGen(args)
	case TaskDebControl:
		debControl(args)
	case TaskDebContents:
		debContents(args)
	case TaskDebContentsDebian:
		debContentsDebian(args)
	case "help", "h", "-help", "--help", "-h":
		log.Printf("Please specify one of: %v", tasks)
		log.Printf("For help on any task, use `debber <task> -h`")
	default:
		log.Printf("Unrecognised task '%s'", task)
		log.Printf("Please specify one of: %v", tasks)
		log.Printf("For help on any task, use `debber <task> -h`")
		os.Exit(1)
	}

}
