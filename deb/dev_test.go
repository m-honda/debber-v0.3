package deb_test

import (
	"github.com/debber/debber-v0.3/deb"
	"log"
	"testing"
)

func Example_buildDevPackage() {


	pkg := deb.NewControl("testpkg", "0.0.2", "me", "me@a", "Dummy package for doing nothing", "testpkg is package ")
	buildFunc := func(dpkg *deb.Control) error {
		// Generate files here.
		return nil
	}
	dpkg := deb.NewDevPackage(pkg)
	err := buildFunc(dpkg)
	if err != nil {
		log.Fatalf("%v", err)
	}
}

func Test_buildDevPackage(t *testing.T) {

	pkg := deb.NewControl("testpkg", "0.0.2", "me", "me@a", "Dummy package for doing nothing", "testpkg is package ")
	buildFunc := func(dpkg *deb.Control) error {
		// Generate files here.
		return nil
	}
	dpkg := deb.NewDevPackage(pkg)
	err := buildFunc(dpkg)
	if err != nil {
		t.Fatalf("%v", err)
	}
}
