package deb

import (
	"path/filepath"
)

const (
	//DebianBinaryVersionDefault is the current version as specified in .deb archives (filename debian-binary)
	DebianBinaryVersionDefault = "2.0"
	//DebianCompatDefault - compatibility. Current version
	DebianCompatDefault = "9"
	//FormatDefault - the format specified in the dsc file (3.0 quilt uses a .debian.gz file rather than a .diff.gz file)
	FormatDefault = "3.0 (quilt)"
	// StatusDefault is unreleased by default. Change this once you're happy with it.
	StatusDefault = "unreleased"

	//SectionDefault - devel seems to be the most common value
	SectionDefault = "devel"
	//PriorityDefault - 'extra' means 'low priority'
	PriorityDefault = "extra"
	//DependsDefault - No dependencies by default
	DependsDefault = ""
	//BuildDependsDefault - debhelper recommended for any package
	BuildDependsDefault = "debhelper (>= 9.1.0)"
	//BuildDependsGoDefault - golang required
	BuildDependsGoDefault = "debhelper (>= 9.1.0), golang-go"

	//StandardsVersionDefault - standards version is specified in the control file
	StandardsVersionDefault = "3.9.4"

	//ArchitectureDefault -'any' is the default architecture for source packages - not for binary debs
	ArchitectureDefault = "any"

	//TemplateDirDefault - the place where control file templates are kept
	TemplateDirDefault = "templates"
	//ResourcesDirDefault - the place where portable files are stored.
	ResourcesDirDefault = "resources"
	//WorkingDirDefault - the directory for build process.
	WorkingDirDefault = "."

	//ExeDirDefault - the default directory for exes within the data archive
	ExeDirDefault                   = "/usr/bin"
	BinaryDataArchiveNameDefault    = "data.tar.gz"
	BinaryControlArchiveNameDefault = "control.tar.gz"

	//OutDirDefault is the default output directory for temp or dist files
	outDirDefault = "target"
	DebianDir     = "debian"
)

const (
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

	BuildDependsFName      = "Build-Depends" // BuildDepends is only required for "sourcedebs".
	BuildDependsIndepFName = "Build-Depends-Indep"
	ConflictsIndepFName    = "Conflicts-Indep"
	BuiltUsingFName        = "Built-Using"

	PriorityFName         = "Priority"
	StandardsVersionFName = "Standards-Version"
	SectionFName          = "Section"
	FormatFName           = "Format"
	StatusFName           = "Status"
	OtherFName            = "Other"
	SourceFName           = "Source"
)

var (
	//TempDirDefault is the default directory for intermediate files
	TempDirDefault = filepath.Join(outDirDefault, "tmp")

	//DistDirDefault is the default directory for built artifacts
	DistDirDefault = outDirDefault

	MaintainerScripts = []string{"postinst", "postrm", "prerm", "preinst"}

	//SourceFields are the fields applicable to Source packages
	//
	// see http://manpages.ubuntu.com/manpages/precise/man5/deb-src-control.5.html://manpages.ubuntu.com/manpages/precise/man5/deb-src-control.5.html
	SourceFields = []string{
		SourceFName,
		MaintainerFName,
		"Uploaders",
		StandardsVersionFName,
		"DM-Upload-Allowed",
		"Homepage",
		"Bugs",
		"Vcs-Arch",
		"Vcs-Bzr",
		"Vcs-Cvs",
		"Vcs-Darcs",
		"Vcs-Git",
		"Vcs-Hg",
		"Vcs-Mtn",
		"Vcs-Svn",
		"Vcs-Browser",
		"Origin",
		SectionFName,
		PriorityFName,
		BuildDependsFName,
		"Build-Depends-Indep",
		"Build-Conflicts",
		"Build-Conflicts-Indep",
	}

	//BinaryFields are the fields applicable to binary packages
	//
	// see http://manpages.ubuntu.com/manpages/precise/man5/deb-src-control.5.html://manpages.ubuntu.com/manpages/precise/man5/deb-src-control.5.html
	BinaryFields = []string{
		PackageFName,
		ArchitectureFName,
		"Package-Type",
		"Subarchitecture",
		"Kernel-Version",
		"Installer-Menu-Item",
		"Essential",
		"Multi-Arch",
		"Tag",
		DescriptionFName,
		DependsFName,
		PreDependsFName,
		RecommendsFName,
		SuggestsFName,
		"Breaks",
		"Enhances",
		"Replaces",
		"Conflicts",
		"Provides",
		"Built-Using",
		PriorityFName,
		SectionFName,
		"Homepage",
	}
)
