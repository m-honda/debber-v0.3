/*
   Copyright 2013 Am Laher

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package deb

import (
	"log"
	"path/filepath"
)

// *DevDebPackage builds a sources-only '-dev' package, which can be used as a BuildDepends dependency. 
// For pure Go packages, this can be cross-platform (architecture == 'all'), but in some cases it might need to be architecture specific
type DevDebPackage struct {
	*Package
	DebFilePath    string
	DestinationGoPathElement string
	BinaryPackage *BinaryPackage
}

// Factory for DevDebPackage
func NewDevPackage(pkg *Package) *DevDebPackage {
	debPath := filepath.Join(pkg.DestDir, pkg.Name+"-dev_"+pkg.Version+".deb")
	return &DevDebPackage{Package: pkg,
		DebFilePath:    debPath,
		DestinationGoPathElement: DEVDEB_GO_PATH_DEFAULT}
}

func (ddpkg *DevDebPackage) InitBinaryPackage() {
	if ddpkg.BinaryPackage == nil {
		//TODO *complete* copy of package. Use reflection??
		devpkg := NewPackage(ddpkg.Name+"-dev", ddpkg.Version, ddpkg.Maintainer)
		devpkg.Description = ddpkg.Description
		devpkg.MaintainerEmail = ddpkg.MaintainerEmail
		devpkg.AdditionalControlData = ddpkg.AdditionalControlData
		devpkg.Architecture = "all"
		devpkg.IsVerbose = ddpkg.IsVerbose
		devpkg.IsRmtemp = ddpkg.IsRmtemp
		ddpkg.BinaryPackage = NewBinaryPackage(devpkg)
	}
	if ddpkg.BinaryPackage.Resources == nil {
		ddpkg.BinaryPackage.Resources = map[string]string{}
	}
}

func (ddpkg *DevDebPackage) BuildWithDefaults() error {
	if ddpkg.BinaryPackage == nil {
		ddpkg.InitBinaryPackage()
	}
	goPathRoot := getGoPathElement(ddpkg.WorkingDir)
	resources, err := globForSources(goPathRoot, ddpkg.WorkingDir, ddpkg.DestinationGoPathElement, []string{ddpkg.TmpDir, ddpkg.DestDir})
	if err != nil {
		return err
	}
	log.Printf("Resources found: %v", resources)
	for k, v := range resources {
		ddpkg.BinaryPackage.Resources[k] = v
	}
	err = ddpkg.BinaryPackage.BuildWithDefaults(Arch_all)
	return err
}


