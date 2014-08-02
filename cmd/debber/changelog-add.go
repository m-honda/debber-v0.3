package main

import (
	"github.com/debber/debber-v0.3/deb"
	"github.com/debber/debber-v0.3/debgen"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

//Adds an entry to the changelog
func changelogAdd(args []string) {
	build := debgen.NewBuildParams()
	fs := InitBuildFlags(cmdName+" "+TaskChangelogAdd, build)
	var version string
	fs.StringVar(&version, "version", "", "Package Version")
	var architecture string
	fs.StringVar(&architecture, "arch", "all", "Architectures [any,386,armhf,amd64,all]")
	var distribution string
	fs.StringVar(&distribution, "distribution", "unstable", "Distribution (unstable is recommended until Debian accept the package into testing/stable)")
	var entry string
	fs.StringVar(&entry, "entry", "", "Changelog entry data")

	err := fs.Parse(os.Args[2:])
	if err != nil {
		log.Fatalf("%v", err)
	}
	if version == "" {
		log.Fatalf("Error: --version is a required flag")
	}
	if entry == "" {
		log.Fatalf("Error: --entry is a required flag")
	}
	controlFilename := filepath.Join(build.DebianDir, "control")
	fi, err := os.Open(controlFilename)
	if os.IsNotExist(err) {
		log.Fatalf("file 'control' not found in debian-dir %s: %v", build.DebianDir, err)
	}
	cfr := deb.NewControlFileReader(fi)
	ctrl, err := cfr.Parse()
	if err != nil {
		log.Fatalf("%v", err)
	}
	(*ctrl)[0].Set("Version", version)
	(*ctrl)[0].Set("Distribution", distribution)
	filename := filepath.Join(build.DebianDir, "changelog")
	templateVars := debgen.NewTemplateData(ctrl)
	templateVars.ChangelogEntry = entry
	err = os.MkdirAll(filepath.Join(build.ResourcesDir, "debian"), 0777)
	if err != nil {
		log.Fatalf("Error making dirs: %v", err)
	}
	_, err = os.Stat(filename)
	if os.IsNotExist(err) {
		tpl, err := template.New("changelog-new").Parse(debgen.TemplateChangelogInitial)
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
		tpl, err := template.New("changelog-add").Parse(debgen.TemplateChangelogAdditionalEntry)
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
