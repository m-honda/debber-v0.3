package debgen_test

import (
	"github.com/debber/debber-v0.3/deb"
	"github.com/debber/debber-v0.3/debgen"
	"log"
)

func Example_genDevPackage() {
	pkg := deb.NewControl("testpkg", "0.0.2", "me", "me@a", "Dummy package for doing nothing", "testpkg is package ")

	ddpkg := deb.NewDevPackage(pkg)
	build := debgen.NewBuildParams()
	build.IsRmtemp = false
	build.Init()
	var err error
	mappedFiles, err := debgen.GlobForGoSources(".", []string{build.TmpDir, build.DestDir})
	if err != nil {
		log.Fatalf("Error building -dev: %v", err)
	}

	err = debgen.GenDevArtifact(ddpkg, build, mappedFiles)
	if err != nil {
		log.Fatalf("Error building -dev: %v", err)
	}

	// Output:
	//
}
