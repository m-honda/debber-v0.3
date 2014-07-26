

package deb

import (
	"path/filepath"
	"os"
	"testing"
)

func TestParseControlFile(t *testing.T) {
	files := []string{ 
		filepath.Join("testdata", "butaca.control"),
		filepath.Join("testdata", "gitso.control"),
		filepath.Join("testdata", "kompas-plugins.control"),
		filepath.Join("testdata", "xkcdMeegoReader.control"),
	}
	for _, filename := range files {
		t.Logf("Package contents of %v:", filename)
		file, err := os.Open(filename)
		if err != nil {
			t.Errorf("cant open file", err)
		}
		cfr := NewControlFileReader(file)
		pkg, err := cfr.Parse()
		if err != nil {
			t.Errorf("cant parse file", err)
		}
		t.Logf("Package contents: %+v", pkg.Paragraphs[0])
	}
}
