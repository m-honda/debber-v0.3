package main

import (
	"fmt"
	"github.com/debber/debber-v0.3/deb"
	"github.com/debber/debber-v0.3/debgen"
	"log"
	"os"
	"path/filepath"
	"text/template"
	"strings"
)

var cmdName = "debber"
/*
func ParseFlags(pkg *deb.Control, params *Params, fs *flag.FlagSet) error {
	err := fs.Parse(os.Args[2:])
	if err != nil {
		return err
	}
	pkg.Paragraphs[0].Set(deb.SourceFName, params.Package)
	pkg.Paragraphs[0].Set(deb.VersionFName, params.Version)
	pkg.Paragraphs[0].Set(deb.MaintainerFName, params.Maintainer)
	pkg.Paragraphs[0].Set(deb.DescriptionFName, params.Description)
	pkg.Paragraphs[0].Set(deb.ArchitectureFName, params.Architecture)
	if len(pkg.Paragraphs) < 2 {
		pkg.Paragraphs = append(pkg.Paragraphs, deb.NewPackage())
	}
	pkg.Paragraphs[1].Set(deb.PackageFName, params.Package)
	if err == nil {
		//validation ...
		err = deb.ValidatePackage(pkg.Paragraphs[0])
		if err != nil {
			println("")
			fmt.Fprintf(os.Stderr, "Usage:\n")
			fs.PrintDefaults()
			println("")
		}
	}
	return err
}
*/
func initDebber(input []string) {
	//set to nothing
	ctrl := deb.NewControlEmpty()
	build := debgen.NewBuildParams()
	build.DestDir = "debian"
	fs := InitBuildFlags(cmdName+" init", build)
	var entry string
	var overwrite bool
	var flavour string
	var pkgName, maintainerName, maintainerEmail, shortDescription, longDescription, architecture string
	fs.StringVar(&pkgName, "name", "", "Package name [required]")
	fs.StringVar(&maintainerName, "maintainer", "", "Maintainer name [required]")
	fs.StringVar(&maintainerEmail, "maintainer-email", "", "Maintainer's email address [required]")
	fs.StringVar(&shortDescription, "desc", "", "Package description [required]")
	fs.StringVar(&longDescription, "long-desc", "", "Package Long description")
	fs.StringVar(&entry, "entry", "Initial project import", "Changelog entry data")
	fs.BoolVar(&overwrite, "overwrite", false, "Overwrite existing files")
	fs.StringVar(&architecture, "architecture", "any", "Package Architecture (any)")
	fs.StringVar(&flavour, "flav", "go:exe", "'flavour' implies a set of defaults - currently, one of 'go:exe', 'go:pkg', 'dev' or ''")
	//TODO flavour
	err := fs.Parse(os.Args[2:])
	if err != nil {
		log.Fatalf("%v", err)
	}
	if pkgName == "" || maintainerName == "" || maintainerEmail == "" || shortDescription == "" || longDescription == ""  {
		log.Fatalf("Required fields: --name, --maintainer, --maintainer-email, --desc, --long-desc")
	}
	//handle spaces in longDescription. TODO: utility function
	longDescriptions := strings.Split(longDescription, "\n")
	longDescription = ""
	for _, ldl := range longDescriptions {
		if longDescription != "" {
			longDescription += "\n"
		}
		longDescription += " " + strings.TrimSpace(ldl)
	}
	(*ctrl)[0].Set(deb.SourceFName, pkgName)
	(*ctrl)[0].Set(deb.MaintainerFName, fmt.Sprintf("%s <%s>", maintainerName, maintainerEmail))
	(*ctrl)[0].Set(deb.DescriptionFName, fmt.Sprintf("%s\n%s", shortDescription, longDescription))
	(*ctrl)[0].Set(deb.ArchitectureFName, architecture)
	if len(*ctrl) < 2 {
		*ctrl = append(*ctrl, deb.NewPackage())
	}
	(*ctrl)[1].Set(deb.PackageFName, pkgName)
	(*ctrl)[1].Set(deb.DescriptionFName, fmt.Sprintf("%s\n%s", shortDescription, longDescription))
	(*ctrl)[1].Set(deb.ArchitectureFName, architecture)
	//-dev package. Optional somehow?
	if len(*ctrl) < 3 {
		*ctrl = append(*ctrl, deb.NewPackage())
	}
	(*ctrl)[2].Set(deb.PackageFName, pkgName+"-dev")
	(*ctrl)[2].Set(deb.ArchitectureFName, "all")
	(*ctrl)[2].Set(deb.DescriptionFName, fmt.Sprintf("%s - development package\n%s", shortDescription, longDescription))
	deb.SetDefaults(ctrl)
	if strings.HasPrefix(flavour, "go:") {
		debgen.ApplyGoDefaults(ctrl)
	}
	if err == nil {
		//validation ...
		err = deb.ValidateControl(ctrl)
		if err != nil {
			println("")
			fmt.Fprintf(os.Stderr, "Usage:\n")
			fs.PrintDefaults()
			println("")
		}
	}
	if err != nil {
		log.Fatalf("%v", err)
	}
	err = build.Init()
	if err != nil {
		log.Fatalf("Error initialising build: %v", err)
	}
	spkg := deb.NewSourcePackage(ctrl)
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
		templateVars := debgen.NewTemplateData(ctrl)
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
		templateVars := debgen.NewTemplateData(ctrl)
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
