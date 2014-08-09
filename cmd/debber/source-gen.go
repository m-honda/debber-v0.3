package main

import (
	"fmt"
	"github.com/debber/debber-v0.3/deb"
	"github.com/debber/debber-v0.3/debgen"
	"os"
	"path/filepath"
)


func sourceGen(input []string) error {
	//set to empty strings because they're being overridden
	//pkg := deb.NewControlEmpty()
	build := debgen.NewBuildParams()
	fs := InitBuildFlags(cmdName+" "+TaskGenSource, build)
	//	fs.StringVar(&pkg.Architecture, "arch", "all", "Architectures [any,386,armhf,amd64,all]")
	fs.StringVar(&build.SourceDir, "sources", ".", "source dir")
	fs.StringVar(&build.SourcesGlob, "sources-glob", debgen.GlobGoSources, "Glob for inclusion of sources")
	fs.StringVar(&build.SourcesRelativeTo, "sources-relative-to", "", "Sources relative to (it will assume relevant gopath element, unless you specify this)")
	fs.StringVar(&build.Version, "version", "", "Package version")

	// parse and validate flags
	err := fs.Parse(os.Args[2:])
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	if build.Version == "" {
		return fmt.Errorf("Error: --version is a required flag")
	}

	if build.SourcesRelativeTo == "" {
		build.SourcesRelativeTo = debgen.GetGoPathElement(build.SourceDir)
	}
	fi, err := os.Open(filepath.Join(build.DebianDir, "control"))
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	cfr := deb.NewControlFileReader(fi)
	ctrl, err := cfr.Parse()
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	spgen, err := debgen.PrepareSourceDebGenerator(ctrl, build)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	err = spgen.GenerateAllDefault()
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	return err
}

