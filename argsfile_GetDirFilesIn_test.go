package argsfile

import (
	"sort"
	"testing"
)

func TestGetDirFilesIn(t *testing.T) {
	dirfiles, err := GetDirFilesIn("testapp", "testdata")

	if err != nil {
		t.Errorf("unexpected error: %v\n", err)
	}

	if dirfiles.DefArgs != "testapp.auto.args" {
		t.Errorf("unexpected Defargs value: %q", dirfiles.DefArgs)
	}

	sort.Strings(dirfiles.ArgsFiles)
	exp := []string{"testapp.de.args", "testapp.en.args", "testapp.es.args", "testapp.fr.args"}
	sort.Strings(exp)
	for _, m := range cmpStrSlc(dirfiles.ArgsFiles, exp) {
		t.Error(m)
	}
}
