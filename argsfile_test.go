package argsfile

import (
	"fmt"
	"testing"
)

func cmpStrSlc(args []string, expArgs []string) (msg []string) {
	if len(args) != len(expArgs) {
		msg = append(msg, fmt.Sprintf("unexpected number of entries: %d != %d", len(args), len(expArgs)))
	} else {
		for i, e := range args {
			if e != expArgs[i] {
				msg = append(msg, fmt.Sprintf("unexpected entry at %d: '%s' != '%s'", i, e, expArgs[i]))
			}
		}
	}

	return msg
}

func TestPullNoAutoArgs(t *testing.T) {
	reduced, found := PullNoAutoArgs([]string{"tool", "--verbose", "--no-auto-args", "file1.txt", "file2.txt"})
	exp := []string{"tool", "--verbose", "file1.txt", "file2.txt"}
	if !found {
		t.Errorf("did not found '--no-auto-args'")
	}
	for _, m := range cmpStrSlc(reduced, exp) {
		t.Error(m)
	}

	reduced, found = PullNoAutoArgs(reduced)
	if found {
		t.Errorf("found unexpected '--no-auto-args' argument")
	}
	for _, m := range cmpStrSlc(reduced, exp) {
		t.Error(m)
	}
}

func TestPullNoArgs(t *testing.T) {
	reduced, found := PullNoArgs([]string{"tool", "--verbose", "--no-args", "file1.txt", "file2.txt", "file3.txt"})
	exp := []string{"tool", "--verbose", "file1.txt", "file2.txt", "file3.txt"}
	if !found {
		t.Errorf("did not found '--no-args'")
	}
	for _, m := range cmpStrSlc(reduced, exp) {
		t.Error(m)
	}

	reduced, found = PullNoAutoArgs(reduced)
	if found {
		t.Errorf("found unexpected '--no-auto-args' argument")
	}
	for _, m := range cmpStrSlc(reduced, exp) {
		t.Error(m)
	}
}
