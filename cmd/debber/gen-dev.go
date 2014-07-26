package main

import (
	"github.com/debber/debber-v0.3/deb"
	"github.com/debber/debber-v0.3/debgen"
	"log"
)

func genDev(input []string) {
	//set to empty strings because they're being overridden
	pkg := deb.NewPackage("", "", "", "")
	build := debgen.NewBuildParams()
	debgen.ApplyGoDefaults(pkg)
	fs, params := InitFlags(cmdName, pkg, build)
//	fs.StringVar(&pkg.Architecture, "arch", "all", "Architectures [any,386,armhf,amd64,all]")
	ddpkg := deb.NewDevPackage(pkg)

	var sourceDir string
	var glob string
	var sourcesRelativeTo string
	var sourcesDestinationDir string
	fs.StringVar(&sourceDir, "sources", build.WorkingDir, "source dir")
	fs.StringVar(&glob, "sources-glob", debgen.GlobGoSources, "Glob for inclusion of sources")
	fs.StringVar(&sourcesRelativeTo, "sources-relative-to", "", "Sources relative to (it will assume relevant gopath element, unless you specify this)")
	fs.StringVar(&sourcesDestinationDir, "sources-destination", debgen.DevGoPathDefault, "Destination dir for sources to be installed")
	err := ParseFlags(pkg, params, fs)
	if err != nil {
		log.Fatalf("%v", err)
	}

	if sourcesRelativeTo == "" {
		sourcesRelativeTo = debgen.GetGoPathElement(sourceDir)
	}
	mappedFiles, err := debgen.GlobForSources(sourcesRelativeTo, sourceDir, glob, sourcesDestinationDir, []string{build.TmpDir, build.DestDir})
	if err != nil {
		log.Fatalf("Error resolving sources: %v", err)
	}
	err = build.Init()
	if err != nil {
		log.Fatalf("Error creating build directories: %v", err)
	}
	err = debgen.GenDevArtifact(ddpkg, build, mappedFiles)
	if err != nil {
		log.Fatalf("Error building -dev: %v", err)
	}

}
