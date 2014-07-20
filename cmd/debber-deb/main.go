package main

import (
	"flag"
	"github.com/debber/debber-v0.3/deb"
	"log"
	"os"
)

func main() {
	name := "debgo-deb"
	log.SetPrefix("[" + name + "] ")

	fs := flag.NewFlagSet(name, flag.ContinueOnError)

	var isControl, isContents, isDebianContents bool
	fs.BoolVar(&isControl, "control", false, "Show control")
	fs.BoolVar(&isContents, "contents", false, "Show contents of data archive")
	fs.BoolVar(&isDebianContents, "debian-contents", false, "Show contents of 'debian' archive (metadata and scripts)")

	//var debFile string
	//fs.StringVar(&debFile, "file", "", ".deb file")
	err := fs.Parse(os.Args[1:])
	if err != nil {
		log.Fatalf("%v", err)
	}
	args := fs.Args()
	if len(args) < 1 {
		log.Fatalf("File not specified")
	}
	if isControl {
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

	} else if isContents {
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
	} else if isDebianContents {
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
	} else {
		log.Fatalf("No command specified")
	}

}
