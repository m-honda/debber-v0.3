package main

import (
	"github.com/laher/debgo-v0.2/deb"
	"github.com/laher/debgo-v0.2/cmd"
	"log"
)

func main() {
	name := "debgo"
	log.SetPrefix("["+name+"] ")
	//set to empty strings because they're being overridden
	pkg := deb.NewGoPackage("","","")

	fs := cmdutils.InitFlags(name, pkg)
	err := cmdutils.ParseFlags(name, pkg, fs)
	if err != nil {
		log.Fatalf("%v", err)
	}
	// TODO determine this platform
	// TODO find executables for this platform
	bpkg := deb.NewBinaryPackage(pkg)
	err = bpkg.BuildAllWithDefaults()
	if err != nil {
		log.Fatalf("%v", err)
	}

}
