package main

import (
	"flag"
	"github.com/debber/debber-v0.3/debgen"
)

func InitFlagsBasic(name string) *flag.FlagSet {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	return fs
}
/*
type Params struct {
	Package      string
	Version      string
	Maintainer   string
	Description  string
	Architecture string
}
*/

func InitBuildFlags(name string, build *debgen.BuildParams) *flag.FlagSet {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.BoolVar(&build.IsRmtemp, "rmtemp", false, "Remove 'temp' dirs")
	fs.BoolVar(&build.IsVerbose, "verbose", false, "Show log messages")
	fs.StringVar(&build.WorkingDir, "working-dir", build.WorkingDir, "Working directory")
	fs.StringVar(&build.TemplateDir, "template-dir", build.TemplateDir, "Template directory")
	fs.StringVar(&build.ResourcesDir, "resources-dir", build.ResourcesDir, "Resources directory")
	fs.StringVar(&build.DebianDir, "debian-dir", build.DebianDir, "'debian' dir (contains control file, changelog, postinst, etc)")
	return fs
}
/*
func InitFlags(name string, pkg *deb.Control, build *debgen.BuildParams) (*flag.FlagSet, *Params) {
	fs := InitBuildFlags(name, build)

	pkgv := new(Params)
	fs.StringVar(&pkgv.Package, "package", "", "Package name")
	fs.StringVar(&pkgv.Version, "version", "", "Package version")
	fs.StringVar(&pkgv.Maintainer, "maintainer", "", "Package maintainer")
	fs.StringVar(&pkgv.Description, "description", "", "Description")
	return fs, pkgv
}
*/

