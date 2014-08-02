package main

import (
	"github.com/debber/debber-v0.3/deb"
	"github.com/debber/debber-v0.3/debgen"
	"log"
	"os"
	"path/filepath"
	"strings"
)
type debGenOpts struct {
	bin386Glob string
	binArmhfGlob string
	binAmd64Glob string
	binAnyGlob string
	resourcesGlob string
	sourcesGlob, sourcesDest string
	version string
	archFilter string
}

func debGen(input []string) {
	build := debgen.NewBuildParams()
	opts := debGenOpts{}
	fs := InitBuildFlags(cmdName+" "+TaskGenDeb, build)
	fs.StringVar(&opts.sourcesGlob, "sources", "**.go", "Glob pattern for sources.")
	fs.StringVar(&opts.bin386Glob, "bin-386", "*386/*", "Glob pattern for binaries for the 386 platform.")
	fs.StringVar(&opts.binArmhfGlob, "bin-armhf", "*armhf/*", "Glob pattern for binaries for the armhf platform.")
	fs.StringVar(&opts.binAmd64Glob, "bin-amd64", "*amd64/*", "Glob pattern for binaries for the amd64 platform.")
	fs.StringVar(&opts.binAnyGlob, "bin-any", "*any/*", "Glob pattern for binaries for *any* platform.")
	fs.StringVar(&opts.sourcesDest, "sources-dest", debgen.DevGoPathDefault + "/src", "directory containing sources.")
	fs.StringVar(&opts.archFilter, "arch-filter", "", "Filter by Architecture. Comma-separated [386,armhf,amd64,all]") //TODO filter outputs by arch?
	fs.StringVar(&opts.resourcesGlob, "resources", "", "directory containing resources for this platform")
	fs.StringVar(&opts.version, "version", "", "Package version")
	err := fs.Parse(os.Args[2:])
	if err != nil {
		log.Fatalf("%v", err)
	}
	if opts.version == "" {
		log.Fatalf("Error: --version is a required flag")
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
	log.Printf("sourcePara: %+v", sourcePara)
	for _, binPara := range ctrl.BinaryParas() {
//TODO check -dev package here ...
		debpara := deb.CopyPara(binPara)
		debpara.Set(deb.VersionFName, opts.version)
		debpara.Set(deb.SectionFName, sourcePara.Get(deb.SectionFName))
		debpara.Set(deb.MaintainerFName, sourcePara.Get(deb.MaintainerFName))
		debpara.Set(deb.PriorityFName, sourcePara.Get(deb.PriorityFName))
		log.Printf("debPara: %+v", debpara)
		sources := []string{}
		//source package. TODO: find a better way to identify a source package.
		if strings.HasSuffix(binPara.Get(deb.PackageFName), "-dev") {
			if opts.sourcesGlob != "" {
				sources, err = filepath.Glob(opts.sourcesGlob)
				if err != nil {
					log.Fatalf("%v", err)
				}
				log.Printf("sources matching glob: %+v", sources)
			}
		
		} else {
			// bin dirs
		}
		//log.Printf("Resources: %v", build.Resources)
		// TODO determine this platform
		//err = bpkg.Build(build, debgen.GenBinaryArtifact)
		artifacts, err := deb.NewWriters(&deb.Control{debpara})
		if err != nil {
			log.Fatalf("%v", err)
		}
		for arch, artifact := range artifacts {
			dgen := debgen.NewDebGenerator(artifact, build)
			for _, source := range sources {
				//NOTE: this should not use filepath.Join because it should always use forward-slash
				dgen.OrigFiles[opts.sourcesDest + "/" + source] = source
			}
			// add resources ...
			err = filepath.Walk(build.ResourcesDir, func(path string, info os.FileInfo, err2 error) error {
				if info != nil && !info.IsDir() {
					rel, err := filepath.Rel(build.ResourcesDir, path)
					if err == nil {
						dgen.OrigFiles[rel] = path
					}
					return err
				}
				return nil
			})
			if err != nil {
				log.Fatalf("%v", err)
			}
			globs := []string{}
			// add binaries ...
			switch arch {
			case "386":
				globs = []string{opts.bin386Glob, opts.binAnyGlob}
			case "amd64":
				globs = []string{opts.binAmd64Glob, opts.binAnyGlob}
			case "armhf":
				globs = []string{opts.binArmhfGlob, opts.binAnyGlob}
			}
			for _, glob := range globs {
				bins, err := filepath.Glob(glob)
				if err != nil {
					log.Fatalf("%v", err)
				}
				log.Printf("Binaries matching glob for '%s': %+v", arch, bins)
				for _, bin := range bins {
					dgen.OrigFiles[deb.ExeDirDefault + "/" + bin] = bin
				}
			}
/*
			archBinDir := filepath.Join(binDir, string(arch))
			err = filepath.Walk(archBinDir, func(path string, info os.FileInfo, err2 error) error {
				if info != nil && !info.IsDir() {
					rel, err := filepath.Rel(binDir, path)
					if err == nil {
						dgen.OrigFiles[rel] = path
					}
					return err
				}
				return nil
			})
*/			
			
			err = dgen.GenerateAllDefault()
			if err != nil {
				log.Fatalf("Error building for '%s': %v", arch, err)
			}
		}
	}
}
