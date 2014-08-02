package main

import (
	"flag"
	"github.com/debber/debber-v0.3/deb"
	"log"
	"os"
)

func debControl(input []string) {
	args := parseFlagsDeb(input)
	for _, debFile := range args {
		rdr, err := os.Open(debFile)
		if err != nil {
			log.Fatalf("%v", err)
		}
		log.Printf("File: %+v", debFile)
		err = deb.DebExtractFileL2(rdr, "control.tar.gz", "control", os.Stdout)
		if err != nil {
			log.Fatalf("%v", err)
		}
	}
}

func debContents(input []string) {
	args := parseFlagsDeb(input)
	for _, debFile := range args {
		rdr, err := os.Open(debFile)
		if err != nil {
			log.Fatalf("%v", err)
		}
		log.Printf("File: %+v", debFile)
		files, err := deb.DebGetContents(rdr, "data.tar.gz")
		if err != nil {
			log.Fatalf("%v", err)
		}
		for _, file := range files {
			log.Printf("%s", file)
		}
	}

}

func debContentsDebian(input []string) {
	args := parseFlagsDeb(input)
	for _, debFile := range args {
		rdr, err := os.Open(debFile)
		if err != nil {
			log.Fatalf("%v", err)
		}
		log.Printf("File: %+v", debFile)
		files, err := deb.DebGetContents(rdr, "control.tar.gz")
		if err != nil {
			log.Fatalf("%v", err)
		}
		for _, file := range files {
			log.Printf("%s", file)
		}
	}

}

func parseFlagsDeb(input []string) []string {
	fs := flag.NewFlagSet(cmdName, flag.ContinueOnError)
	err := fs.Parse(input)
	if err != nil {
		log.Fatalf("%v", err)
	}
	args := fs.Args()
	if len(args) < 1 {
		log.Fatalf("File not specified")
	}
	return args
}
