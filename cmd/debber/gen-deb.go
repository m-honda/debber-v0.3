package main

import (
	"github.com/debber/debber-v0.3/cmd"
	"github.com/debber/debber-v0.3/deb"
	"github.com/debber/debber-v0.3/debgen"
	"log"
	"os"
	"path/filepath"
)

func genDeb(input []string) {
	build := debgen.NewBuildParams()
	fs := cmdutils.InitBuildFlags(cmdName+" "+TaskGenDeb, build)
	var binDir string
	//var resourcesDir string
	var arch, version string
	fs.StringVar(&binDir, "binaries", "", "directory containing binaries for each architecture. Directory names should end with the architecture")
	fs.StringVar(&arch, "arch", "any", "Architectures [any,386,armhf,amd64,all]")
	//fs.StringVar(&resourcesDir, "resources", "", "directory containing resources for this platform")
	fs.StringVar(&version, "version", "", "Package version")
	err := fs.Parse(os.Args[2:])
	if err != nil {
		log.Fatalf("%v", err)
	}
	if version == "" {
		log.Fatalf("-version is required", version)
	}
	fi, err := os.Open(filepath.Join(build.DebianDir, "control"))
	if err != nil {
		log.Fatalf("%v", err)
	}
	cfr := deb.NewControlFileReader(fi)
	pkg, err := cfr.Parse()
	pkg.Paragraphs[0].Set(deb.VersionFName, version)
	if err != nil {
		log.Fatalf("%v", err)
	}


	if err != nil {
		log.Fatalf("%v", err)
	}
	err = build.Init()
	if err != nil {
		log.Fatalf("%v", err)
	}

	//log.Printf("Resources: %v", build.Resources)
	// TODO determine this platform
	//err = bpkg.Build(build, debgen.GenBinaryArtifact)
	artifacts, err := deb.NewDebWriters(pkg)
	if err != nil {
		log.Fatalf("%v", err)
	}
	for arch, artifact := range artifacts {
		dgen := debgen.NewDebGenerator(artifact, build)
		err = filepath.Walk(build.ResourcesDir, func(path string, info os.FileInfo, err2 error) error {
			if info != nil && !info.IsDir() {
				rel, err := filepath.Rel(resourcesDir, path)
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
		err = dgen.GenerateAllDefault()
		if err != nil {
			log.Fatalf("Error building for '%s': %v", arch, err)
		}
	}
}
