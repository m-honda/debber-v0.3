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
	"strings"
)

// Control is the base unit for this library.
// A *Control contains one or more paragraphs
type Control []*Package


// NewEmptyControl returns a package with one empty paragraph and an empty map of ExtraData
func NewControlEmpty() *Control {
	ctrl := &Control{NewPackage()}
	return ctrl
}

// NewControl is a factory for a Control. Name, Version, Maintainer and Description are mandatory.
func NewControlDefault(name, maintainerName, maintainerEmail, shortDescription, longDescription string, addDevPackage bool) *Control {
	ctrl := NewControlEmpty()
	(*ctrl)[0].Set(SourceFName, name)
	(*ctrl)[0].Set(MaintainerFName, fmt.Sprintf("%s <%s>", maintainerName, maintainerEmail))
	(*ctrl)[0].Set(DescriptionFName, fmt.Sprintf("%s\n%s", shortDescription, longDescription))
	//BuildDepends is empty...
	//add binary package
	*ctrl = append(*ctrl, NewPackage())
	(*ctrl)[1].Set(PackageFName, name)
	(*ctrl)[1].Set(DescriptionFName, fmt.Sprintf("%s\n%s", shortDescription, longDescription))
	//depends is empty
	if addDevPackage {
		*ctrl = append(*ctrl, NewPackage())
		(*ctrl)[2].Set(PackageFName, name+"-dev")
		(*ctrl)[2].Set(ArchitectureFName, "all")
		(*ctrl)[2].Set(DescriptionFName, fmt.Sprintf("%s - development package\n%s", shortDescription, longDescription))
	}
	SetDefaults(ctrl)
	return ctrl
}

func (ctrl *Control) NewDevPackage() *Package {
	name := ctrl.Get(PackageFName)
	desc := ctrl.Get(DescriptionFName)
	devPara := NewPackage()
	devPara.Set(PackageFName, name+"-dev")
	devPara.Set(ArchitectureFName, "all")
	sp := strings.SplitN(desc, "\n", 2)
	shortDescription := sp[0]
	longDescription := ""
	if len(sp) > 1 {
		longDescription = sp[1]
	}
	devPara.Set(DescriptionFName, fmt.Sprintf("%s - development package\n%s", shortDescription, longDescription))
	return devPara
}

// Sets fields which can be initialised appropriately
// note that Source and Binary packages are detected by the presence of a Source or Package field, respectively.
func SetDefaults(ctrl *Control) {

	for _, pkg := range *ctrl {
		if pkg.Get(SourceFName) != "" {
			pkg.Set(PriorityFName, PriorityDefault)
			pkg.Set(StandardsVersionFName, StandardsVersionDefault)
			pkg.Set(SectionFName, SectionDefault)
			pkg.Set(FormatFName, FormatDefault)
			pkg.Set(StatusFName, StatusDefault)
		}
	}
	//ctrl.MappedFiles = map[string]string{}
	for _, pkg := range *ctrl {
		if pkg.Get(PackageFName) != "" {
			if pkg.Get(ArchitectureFName) == "" {
				pkg.Set(ArchitectureFName, "any") //default ...
			}
		}
	}
}

// GetArches resolves architecture(s) and return as a slice
func (ctrl *Control) GetArches() ([]Architecture, error) {
	_, arch, exists := (*ctrl)[0].GetExtended(ArchitectureFName)
	if exists {
		arches, err := resolveArches(arch)
		return arches, err
	}
	if len(*ctrl) > 1 {
		_, arch2, exists2 := (*ctrl)[1].GetExtended(ArchitectureFName)
		if exists2 {
			arches, err := resolveArches(arch2)
			return arches, err
		}
	}
	return nil, fmt.Errorf("Architecture field not set")

}

//Get finds the first occurence of the specified value, checking each paragraph in turn
func (ctrl *Control) Get(key string) string {
	for _, paragraph := range *ctrl {
		//return (*ctrl)[0].Get(key)
		_, val, exists := paragraph.GetExtended(key)
		if exists {
			return val
		}
	}
	//not found
	return ""
}

// Copy all fields
func Copy(ctrl *Control) *Control {
	nctrl := NewControlEmpty()
	for _, para := range *ctrl {
		*nctrl = append(*nctrl, CopyPara(para))
	}
	return nctrl
}

func (ctrl *Control) BinaryParas() []*Package {
	paras := []*Package{}
	nkey := PackageFName
	for _, para := range *ctrl {
		v := para.Get(nkey)
		if v != "" {
			paras = append(paras, para)
		}
	}
	return paras
}

func (ctrl *Control) SourceParas() []*Package {
	paras := []*Package{}
	nkey := SourceFName
	for _, para := range *ctrl {
		v := para.Get(nkey)
		if v != "" {
			paras = append(paras, para)
		}
	}
	return paras
}

func (ctrl *Control) GetParasByField(key string, val string) []*Package {
	paras := []*Package{}
	nkey := NormaliseFieldKey(key)
	for _, para := range *ctrl {
		v := para.Get(nkey)
		if val == v {
			paras = append(paras, para)
		}
	}
	return paras
}
