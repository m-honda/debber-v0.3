package main

import (
	"flag"
	"fmt"
	"github.com/debber/debber-v0.3/deb"
	"github.com/debber/debber-v0.3/debgen"
	"os"
)

func InitFlagsBasic(name string) *flag.FlagSet {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	return fs
}

type Params struct {
	Package      string
	Version      string
	Maintainer   string
	Description  string
	Architecture string
}

func InitBuildFlags(name string, build *debgen.BuildParams) (*flag.FlagSet) {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.BoolVar(&build.IsRmtemp, "rmtemp", false, "Remove 'temp' dirs")
	fs.BoolVar(&build.IsVerbose, "verbose", false, "Show log messages")
	fs.StringVar(&build.WorkingDir, "working-dir", build.WorkingDir, "Working directory")
	fs.StringVar(&build.TemplateDir, "template-dir", build.TemplateDir, "Template directory")
	fs.StringVar(&build.ResourcesDir, "resources-dir", build.ResourcesDir, "Resources directory")
	fs.StringVar(&build.DebianDir, "debian-dir", build.DebianDir, "'debian' dir (contains control file, changelog, etc)")
	return fs
}

func InitFlags(name string, pkg *deb.Package, build *debgen.BuildParams) (*flag.FlagSet, *Params) {
	fs := InitBuildFlags(name, build)

	pkgv := new(Params)
	fs.StringVar(&pkgv.Package, "package", "", "Package name")
	fs.StringVar(&pkgv.Version, "version", "", "Package version")
	fs.StringVar(&pkgv.Maintainer, "maintainer", "", "Package maintainer")
	fs.StringVar(&pkgv.Description, "description", "", "Description")
	return fs, pkgv
}

func ParseFlags(pkg *deb.Package, params *Params, fs *flag.FlagSet) error {
	err := fs.Parse(os.Args[2:])
	if err != nil {
		return err
	}
	pkg.Paragraphs[0].Set(deb.PackageFName, params.Package)
	pkg.Paragraphs[0].Set(deb.VersionFName, params.Version)
	pkg.Paragraphs[0].Set(deb.MaintainerFName, params.Maintainer)
	pkg.Paragraphs[0].Set(deb.DescriptionFName, params.Description)
	pkg.Paragraphs[0].Set(deb.ArchitectureFName, params.Architecture)
	if err == nil {
		//validation ...
		err = deb.ValidatePackage(pkg)
		if err != nil {
			println("")
			fmt.Fprintf(os.Stderr, "Usage:\n")
			fs.PrintDefaults()
			println("")
		}
	}
	return err
}
