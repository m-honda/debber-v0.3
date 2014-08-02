package main

import (
	"github.com/debber/debber-v0.3/deb"
	"github.com/debber/debber-v0.3/debgen"
	"log"
	"os"
	"path/filepath"
)

type sourceGenOpts struct {
	sourceDir string
	sourcesRelativeTo string
	sourcesGlob string
	version string
}

func sourceGen(input []string) {
	//set to empty strings because they're being overridden
	//pkg := deb.NewControlEmpty()
	build := debgen.NewBuildParams()
	opts := sourceGenOpts{}
	fs := InitBuildFlags(cmdName+" "+TaskGenSource, build)
	//	fs.StringVar(&pkg.Architecture, "arch", "all", "Architectures [any,386,armhf,amd64,all]")
	fs.StringVar(&opts.sourceDir, "sources", ".", "source dir")
	fs.StringVar(&opts.sourcesGlob, "sources-glob", debgen.GlobGoSources, "Glob for inclusion of sources")
	fs.StringVar(&opts.sourcesRelativeTo, "sources-relative-to", "", "Sources relative to (it will assume relevant gopath element, unless you specify this)")
	fs.StringVar(&opts.version, "version", "", "Package version")

	// parse and validate flags
	err := fs.Parse(os.Args[2:])
	if err != nil {
		log.Fatalf("%v", err)
	}
	if opts.version == "" {
		log.Fatalf("Error: --version is a required flag")
	}


	if opts.sourcesRelativeTo == "" {
		opts.sourcesRelativeTo = debgen.GetGoPathElement(opts.sourceDir)
	}

	fi, err := os.Open(filepath.Join(build.DebianDir, "control"))
	if err != nil {
		log.Fatalf("%v", err)
	}

	cfr := deb.NewControlFileReader(fi)
	ctrl, err := cfr.Parse()
	if err != nil {
		log.Fatalf("%v", err)
	}

	err = build.Init()
	if err != nil {
		log.Fatalf("%v", err)
	}
	sourcePara := ctrl.SourceParas()[0]
	//sourcePara := deb.CopyPara(sp)
	sourcePara.Set(deb.VersionFName, opts.version)
	sourcePara.Set(deb.FormatFName, deb.FormatDefault)

	//Build ...
	spkg := deb.NewSourcePackage(ctrl)
	sourcesDestinationDir := sourcePara.Get(deb.SourceFName) + "_" + sourcePara.Get(deb.VersionFName)
	spgen := debgen.NewSourcePackageGenerator(spkg, build)
	spgen.OrigFiles, err = debgen.GlobForSources(opts.sourcesRelativeTo, opts.sourceDir, opts.sourcesGlob, sourcesDestinationDir, []string{build.TmpDir, build.DestDir})
	if err != nil {
		log.Fatalf("Error resolving sources: %v", err)
	}
	err = spgen.GenerateAllDefault()
	if err != nil {
		log.Fatalf("%v", err)
	}

}
