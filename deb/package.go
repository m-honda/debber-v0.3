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
	//	"log"
	//	"reflect"
	"fmt"
	"strings"
)

// Package is the base unit for this library.
// A *Package contains metadata.
type Package struct {
	Paragraphs []*Paragraph
	ExtraData  map[string]interface{} // Optional for templates
}

type Paragraph struct {
	controlData           map[string]string // This map should contain all keys/values, with keys using the standard camel case
	controlDataCaseLookup map[string]string // Case sensitive lookups. TODO
}

var (
	PackageFName     = "Package"
	VersionFName     = "Version"
	DescriptionFName = "Description"
	MaintainerFName  = "Maintainer"

	ArchitectureFName = "Architecture" // Supported values: "all", "x386", "amd64", "armhf". TODO: armel

	DependsFName    = "Depends" // Depends
	RecommendsFName = "Recommends"
	SuggestsFName   = "Suggests"
	EnhancesFName   = "Enhances"
	PreDependsFName = "PreDepends"
	ConflictsFName  = "Conflicts"
	BreaksFName     = "Breaks"
	ProvidesFName   = "Provides"
	ReplacesFName   = "Replaces"

	BuildDependsFName      = "BuildDepends" // BuildDepends is only required for "sourcedebs".
	BuildDependsIndepFName = "BuildDependsIndep"
	ConflictsIndepFName    = "ConflictsIndep"
	BuiltUsingFName        = "BuiltUsing"

	PriorityFName         = "Priority"
	StandardsVersionFName = "StandardsVersion"
	SectionFName          = "Section"
	FormatFName           = "Format"
	StatusFName           = "Status"
	OtherFName            = "Other"
	SourceFName           = "Source"
)

func newParagraph() *Paragraph {
	para := &Paragraph{controlData: map[string]string{}, controlDataCaseLookup: map[string]string{}}
	return para
}

// NewEmptyPackage returns a package with one empty paragraph and an empty map of ExtraData
func NewEmptyPackage() *Package {
	pkg := &Package{Paragraphs: []*Paragraph{newParagraph()}, ExtraData: map[string]interface{}{}}
	return pkg
}



// NewPackage is a factory for a Package. Name, Version, Maintainer and Description are mandatory.
func NewPackage(name, version, maintainer, description string) *Package {
	pkg := NewEmptyPackage()
	pkg.Paragraphs[0].Set(PackageFName, name)
	pkg.Paragraphs[0].Set(VersionFName, version)
	pkg.Paragraphs[0].Set(MaintainerFName, maintainer)
	pkg.Paragraphs[0].Set(DescriptionFName, description)
	SetDefaults(pkg)
	return pkg
}

// Sets fields which can be initialised appropriately
func SetDefaults(pkg *Package) {
	pkg.Paragraphs[0].Set(ArchitectureFName, "any") //default ...
	pkg.Paragraphs[0].Set(PriorityFName, PriorityDefault)
	pkg.Paragraphs[0].Set(StandardsVersionFName, StandardsVersionDefault)
	pkg.Paragraphs[0].Set(SectionFName, SectionDefault)
	pkg.Paragraphs[0].Set(FormatFName, FormatDefault)
	pkg.Paragraphs[0].Set(StatusFName, StatusDefault)
	//pkg.MappedFiles = map[string]string{}
}

// GetArches resolves architecture(s) and return as a slice
func (pkg *Package) GetArches() ([]Architecture, error) {
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

// Set sets a control field by name
func (pkg *Paragraph) Set(key, value string) error {
	existingKey, exists := pkg.controlDataCaseLookup[strings.ToLower(key)]
	if exists && existingKey != key {
		return fmt.Errorf("Key exists with different case. %s != %s", key, existingKey)
	}
	pkg.controlData[key] = value
	pkg.controlDataCaseLookup[strings.ToLower(key)] = key
	return nil
}

//Get finds the first occurence of the specified value, checking each paragraph in turn
func (pkg *Package) Get(key string) string {
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

// GetExtended gets a control field by name, returning key, value & 'exists'
func (pkg *Paragraph) GetExtended(key string) (string, string, bool) {
	myKey, exists := pkg.controlDataCaseLookup[strings.ToLower(key)]
	if !exists {
		return "", "", exists
	}
	val, exists := pkg.controlData[myKey]
	return myKey, val, exists
}

func (pkg *Paragraph) Get(key string) string {
	myKey, exists := pkg.controlDataCaseLookup[strings.ToLower(key)]
	if !exists {
		return ""
	}
	val, exists := pkg.controlData[myKey]
	return val
}

// Copy all fields
func Copy(pkg *Package) *Package {
	npkg := NewEmptyPackage()
	npkg.Paragraphs = []*Paragraph{}
	for _, para := range pkg.Paragraphs {
		npkg.Paragraphs = append(npkg.Paragraphs, CopyPara(para))
	}
	return npkg
}
func CopyPara(pkg *Paragraph) *Paragraph {
	npkg := newParagraph()
	for k, v := range pkg.controlData {
		npkg.controlData[k] = v
	}
	for k, v := range pkg.controlDataCaseLookup {
		npkg.controlDataCaseLookup[k] = v
	}
	return npkg
}

/*
func Copy(pkg *Package) *Package {
	//ptype := reflect.TypeOf(pkg)
	npkg := &Package{}
	pkgVal := reflect.ValueOf(pkg).Elem()
	npkgVal := reflect.ValueOf(npkg).Elem()
	ptype := pkgVal.Type()
	for i := 0; i < ptype.NumField(); i++ {
		source := pkgVal.Field(i)
		dest := npkgVal.Field(i)
		log.Printf("%v => %v", source, dest)
		dest.Set(source)
	}
	return npkg
}
*/
