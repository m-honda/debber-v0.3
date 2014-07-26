package main

import (
	"github.com/debber/debber-v0.3/deb"
	"github.com/debber/debber-v0.3/debgen"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

var cmdName = "debber"

func initDebber(input []string) {
	//set to nothing
	pkg := deb.NewEmptyControl()
	build := debgen.NewBuildParams()
	build.DestDir = "debian"
	fs, params := InitFlags(cmdName, pkg, build)
	var entry string
	var overwrite bool
	var flavour string
	fs.StringVar(&entry, "entry", "Initial project import", "Changelog entry data")
	fs.BoolVar(&overwrite, "overwrite", false, "Overwrite existing files")
	fs.StringVar(&params.Architecture, "architecture", "any", "Package Architecture (any)")
	fs.StringVar(&flavour, "flav", "go:exe", "'flavour' implies a set of defaults - currently, one of 'go:exe', 'go:pkg', 'dev' or ''")
	//TODO flavour
	err := ParseFlags(pkg, params, fs)
	if err != nil {
		log.Fatalf("%v", err)
	}
	if flavour == "go:exe" {
		debgen.ApplyGoDefaults(pkg)
	}
	err = build.Init()
	if err != nil {
		log.Fatalf("Error initialising build: %v", err)
	}
	spkg := deb.NewSourcePackage(pkg)
	spgen := debgen.NewSourcePackageGenerator(spkg, build)
	if flavour == "go:exe" {
		spgen.ApplyDefaultsPureGo()
	}

	//create control file
	filename := filepath.Join(build.DebianDir, "control")
	_, err = os.Stat(filename)
	if os.IsNotExist(err) || overwrite {
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
	if os.IsNotExist(err) || overwrite {
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
	if os.IsNotExist(err) || overwrite {
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

