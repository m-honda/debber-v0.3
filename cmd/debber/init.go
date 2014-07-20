package main

import (
	"flag"
	"github.com/debber/debber-v0.3/cmd"
	"github.com/debber/debber-v0.3/deb"
	"github.com/debber/debber-v0.3/debgen"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"text/template"
)
var cmdName = "debber"
func initDebber(input []string) {
	pkg := deb.NewPackage("", "", "", "")
	build := debgen.NewBuildParams()
	build.DestDir = "debian"
	debgen.ApplyGoDefaults(pkg)
	fs := cmdutils.InitFlags(cmdName, pkg, build)
	var entry string
	fs.StringVar(&entry, "entry", "Initial project import", "Changelog entry data")
	err := fs.Parse(input)
	if err != nil {
		log.Fatalf("%v", err)
	}
	err = deb.ValidatePackage(pkg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Usage of %s\n", cmdName)
		fs.PrintDefaults()
		log.Fatalf("%v", err)
	}

	err = build.Init()
	if err != nil {
		log.Fatalf("Error initialising build: %v", err)
	}

	//create DSC
	//create changelog
	spkg := deb.NewSourcePackage(pkg)
	spgen := debgen.NewSourcePackageGenerator(spkg, build) 
	spgen.ApplyDefaultsPureGo()

	filename := filepath.Join(build.DebianDir, "control")
	_, err = os.Stat(filename)
	if os.IsNotExist(err) {
		err = spgen.GenSourceControlFile()
		if err != nil {
			log.Fatalf("Error generating control file: %v", err)
		}
	} else if err != nil {
		log.Fatalf("Error generating control file: %v", err)
	} else {
		log.Printf("%s already exists", filename)
	}
	//changelog file
	filename = filepath.Join(build.DebianDir, "changelog")
	_, err = os.Stat(filename)
	if os.IsNotExist(err) {
		templateVars := debgen.NewTemplateData(pkg)
		templateVars.ChangelogEntry = entry
		tpl, err := template.New("template").Parse(debgen.TemplateChangelogInitial)
		if err != nil {
			log.Fatalf("Error parsing template: %v", err)
		}
		//create ..
		f, err := os.Create(filename)
		if err != nil {
			log.Fatalf("Error creating file: %v", err)
		}
		defer f.Close()
		err = tpl.Execute(f, templateVars)
		if err != nil {
			log.Fatalf("Error executing template: %v", err)
		}
		err = f.Close()
		if err != nil {
			log.Fatalf("Error closing written file: %v", err)
		}
	} else if err != nil {
		log.Fatalf("Error reading existing changelog: %v", err)
	} else {
		log.Printf("%s already exists", filename)
	}
	//rules file ...
	filename = filepath.Join(build.DebianDir, "rules")
	_, err = os.Stat(filename)
	if os.IsNotExist(err) {
		templateVars := debgen.NewTemplateData(pkg)
		tpl, err := template.New("template").Parse(spgen.TemplateStrings["rules"])
		if err != nil {
			log.Fatalf("Error parsing template: %v", err)
		}
		//create ..
		f, err := os.Create(filename)
		if err != nil {
			log.Fatalf("Error creating file: %v", err)
		}
		defer f.Close()
		err = tpl.Execute(f, templateVars)
		if err != nil {
			log.Fatalf("Error executing template: %v", err)
		}
		err = f.Close()
		if err != nil {
			log.Fatalf("Error closing written file: %v", err)
		}
	} else if err != nil {
		log.Fatalf("Error reading existing rules file: %v", err)
	} else {
		log.Printf("%s already exists", filename)
	}
}

func genChangelog(args []string) {
	pkg := deb.NewPackage("", "", "", "")
	build := debgen.NewBuildParams()
	debgen.ApplyGoDefaults(pkg)
	fs := cmdutils.InitFlags(cmdName, pkg, build)
	fs.StringVar(&pkg.Architecture, "arch", "all", "Architectures [any,386,armhf,amd64,all]")
	var entry string
	fs.StringVar(&entry, "entry", "", "Changelog entry data")

	err := cmdutils.ParseFlags(cmdName, pkg, fs)
	if err != nil {
		log.Fatalf("%v", err)
	}
	if entry == "" {
		log.Fatalf("Error: --entry is a required flag")

	}
	filename := filepath.Join(build.DebianDir, "changelog")
	templateVars := debgen.NewTemplateData(pkg)
	templateVars.ChangelogEntry = entry
	err = os.MkdirAll(filepath.Join(build.ResourcesDir, "debian"), 0777)
	if err != nil {
		log.Fatalf("Error making dirs: %v", err)
	}

	_, err = os.Stat(filename)
	if os.IsNotExist(err) {
		tpl, err := template.New("template").Parse(debgen.TemplateChangelogInitial)
		if err != nil {
			log.Fatalf("Error parsing template: %v", err)
		}
		//create ..
		f, err := os.Create(filename)
		if err != nil {
			log.Fatalf("Error creating file: %v", err)
		}
		defer f.Close()
		err = tpl.Execute(f, templateVars)
		if err != nil {
			log.Fatalf("Error executing template: %v", err)
		}
		err = f.Close()
		if err != nil {
			log.Fatalf("Error closing written file: %v", err)
		}
	} else if err != nil {
		log.Fatalf("Error reading existing changelog: %v", err)
	} else {
		tpl, err := template.New("template").Parse(debgen.TemplateChangelogAdditionalEntry)
		if err != nil {
			log.Fatalf("Error parsing template: %v", err)
		}
		//append..
		f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatalf("Error opening file: %v", err)
		}
		defer f.Close()
		err = tpl.Execute(f, templateVars)
		if err != nil {
			log.Fatalf("Error executing template: %v", err)
		}
		err = f.Close()
		if err != nil {
			log.Fatalf("Error closing written file: %v", err)
		}
	}

}

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


