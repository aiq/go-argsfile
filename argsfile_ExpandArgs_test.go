package argsfile

import (
	"testing"
)

func TestExpandArgs(t *testing.T) {
	{
		args, err := ExpandArgs([]string{"values", "before", "--args", "testdata/expand.args", "values", "after"})
		if err != nil {
			t.Errorf("unexpected error: %v\n", err)
		}

		exp := []string{"values", "before", "-i", "1", "second value", "third value", "foo=bar", "values", "after"}
		for _, m := range cmpStrSlc(args, exp) {
			t.Error(m)
		}
	}

	//**************************************************************************
	{
		args, err := ExpandArgs([]string{"testapp", "--args", "testdata/testapp.auto.args", "tail"})
		if err != nil {
			t.Errorf("unexpected error: %v\n", err)
		}

		exp := []string{"testapp", "auto", "args", "values", "tail"}
		for _, m := range cmpStrSlc(args, exp) {
			t.Error(m)
		}
	}

	//**************************************************************************
	{
		args, err := ExpandArgs([]string{"testapp", "--args", "testdata/testapp.auto.args"})
		if err != nil {
			t.Errorf("unexpected error: %v\n", err)
		}

		exp := []string{"testapp", "auto", "args", "values"}
		for _, m := range cmpStrSlc(args, exp) {
			t.Error(m)
		}
	}
}
