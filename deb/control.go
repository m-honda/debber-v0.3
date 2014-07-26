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
	"fmt"
)

// Control is the base unit for this library.
// A *Control contains one or more paragraphs
type Control struct {
	Paragraphs []*Package
	ExtraData  map[string]interface{} // Optional for templates
}


// NewEmptyControl returns a package with one empty paragraph and an empty map of ExtraData
func NewEmptyControl() *Control {
	pkg := &Control{Paragraphs: []*Package{NewPackage()}, ExtraData: map[string]interface{}{}}
	return pkg
}

// NewControl is a factory for a Control. Name, Version, Maintainer and Description are mandatory.
func NewControl(name, version, maintainerName, maintainerEmail, shortDescription, longDescription string) *Control {
	pkg := NewEmptyControl()
	pkg.Paragraphs[0].Set(SourceFName, name)
	pkg.Paragraphs[0].Set(VersionFName, version)
	pkg.Paragraphs[0].Set(MaintainerFName, fmt.Sprintf("%s <%s>", maintainerName, maintainerEmail))
	pkg.Paragraphs[0].Set(DescriptionFName, fmt.Sprintf("%s\n%s", shortDescription, longDescription))
	pkg.Paragraphs = append(pkg.Paragraphs, NewPackage())
	SetDefaults(pkg)
	return pkg
}

// Sets fields which can be initialised appropriately
func SetDefaults(pkg *Control) {
	pkg.Paragraphs[0].Set(PriorityFName, PriorityDefault)
	pkg.Paragraphs[0].Set(StandardsVersionFName, StandardsVersionDefault)
	pkg.Paragraphs[0].Set(SectionFName, SectionDefault)
	pkg.Paragraphs[0].Set(FormatFName, FormatDefault)
	pkg.Paragraphs[0].Set(StatusFName, StatusDefault)
	//pkg.MappedFiles = map[string]string{}
	if len(pkg.Paragraphs) > 1 {
		pkg.Paragraphs[1].Set(ArchitectureFName, "any") //default ...
	}
}

// GetArches resolves architecture(s) and return as a slice
func (pkg *Control) GetArches() ([]Architecture, error) {
	_, arch, exists := pkg.Paragraphs[0].GetExtended(ArchitectureFName)
	if exists {
		arches, err := resolveArches(arch)
		return arches, err
	}
	if len(pkg.Paragraphs) > 1 {
		_, arch2, exists2 := pkg.Paragraphs[1].GetExtended(ArchitectureFName)
		if exists2 {
			arches, err := resolveArches(arch2)
			return arches, err
		}
	}
	return nil, fmt.Errorf("Architecture field not set")

}

//Get finds the first occurence of the specified value, checking each paragraph in turn
func (pkg *Control) Get(key string) string {
	for _, paragraph := range pkg.Paragraphs {
		//return pkg.Paragraphs[0].Get(key)
		_, val, exists := paragraph.GetExtended(key)
		if exists {
			return val
		}
	}
	//not found
	return ""
}

// Copy all fields
func Copy(pkg *Control) *Control {
	npkg := NewEmptyControl()
	npkg.Paragraphs = []*Package{}
	for _, para := range pkg.Paragraphs {
		npkg.Paragraphs = append(npkg.Paragraphs, CopyPara(para))
	}
	return npkg
}


