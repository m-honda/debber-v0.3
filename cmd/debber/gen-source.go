package main

import (
	"github.com/debber/debber-v0.3/deb"
	"github.com/debber/debber-v0.3/debgen"
	"log"
)

func genSource(input []string) {
	//set to empty strings because they're being overridden
	pkg := deb.NewEmptyControl()
	build := debgen.NewBuildParams()
	debgen.ApplyGoDefaults(pkg)
	fs, params := InitFlags(cmdName, pkg, build)
//	fs.StringVar(&pkg.Architecture, "arch", "all", "Architectures [any,386,armhf,amd64,all]")

	var sourceDir string
	var glob string
	var sourcesRelativeTo string
	fs.StringVar(&sourceDir, "sources", ".", "source dir")
	fs.StringVar(&glob, "sources-glob", debgen.GlobGoSources, "Glob for inclusion of sources")
	fs.StringVar(&sourcesRelativeTo, "sources-relative-to", "", "Sources relative to (it will assume relevant gopath element, unless you specify this)")
	err := ParseFlags(pkg, params, fs)
	if err != nil {
		log.Fatalf("%v", err)
	}
	err = build.Init()
	if err != nil {
		log.Fatalf("Error creating build directories: %v", err)
	}
	if sourcesRelativeTo == "" {
		sourcesRelativeTo = debgen.GetGoPathElement(sourceDir)
	}
	spkg := deb.NewSourcePackage(pkg)
	sourcesDestinationDir := pkg.Get(deb.PackageFName) + "_" + pkg.Get(deb.VersionFName)
	spgen := debgen.NewSourcePackageGenerator(spkg, build)
	spgen.OrigFiles, err = debgen.GlobForSources(sourcesRelativeTo, sourceDir, glob, sourcesDestinationDir, []string{build.TmpDir, build.DestDir})
	if err != nil {
		log.Fatalf("Error resolving sources: %v", err)
	}
	err = spgen.GenerateAllDefault()
	if err != nil {
		log.Fatalf("%v", err)
	}

}
