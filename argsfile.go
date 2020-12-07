// Package argsfile supports the handling of args files.
//
// See https://args.aiq.dk/ for more information about the args file format.
package argsfile

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
)

func appendToLast(slc []string, conc string, tail string) []string {
	pos := len(slc) - 1
	last := slc[pos]
	slc[pos] = last + conc + tail
	return slc
}

func isWhiteSpace(c rune) bool {
	return c == ' ' || c == '\t' || c == '\r' || c == '\n'
}

func parseShellLine(line string) []string {
	args := []string{}

	builder := &strings.Builder{}
	escaped := false
	doubleQ := false
	singleQ := false
	backQ := false

	for _, r := range line {
		if escaped {
			builder.WriteRune(r)
			escaped = false
		} else if r == '\\' {
			if singleQ {
				builder.WriteRune(r)
			} else {
				escaped = true
			}
		} else if isWhiteSpace(r) {
			if singleQ || doubleQ || backQ {
				builder.WriteRune(r)
			} else {
				args = append(args, builder.String())
				builder = &strings.Builder{}
			}
		} else if r == '`' {
			if singleQ || doubleQ {
				builder.WriteRune(r)
			} else {
				backQ = !backQ
			}
		} else if r == '"' {
			if singleQ || backQ {
				builder.WriteRune(r)
			} else {
				doubleQ = !doubleQ
			}
		} else if r == '\'' {
			if doubleQ || backQ {
				builder.WriteRune(r)
			} else {
				singleQ = !singleQ
			}
		} else {
			builder.WriteRune(r)
		}
	}

	if builder.Len() > 0 {
		args = append(args, builder.String())
	}

	return args
}

func stringIndex(slc []string, val string) (idx int, ok bool) {
	for i, str := range slc {
		if str == val {
			return i, true
		}
	}

	return idx, false
}

func selectArgsFile(files []string) (f string, err error) {
	fmt.Printf("args files in the working directory:\n\n")
	for i, f := range files {
		fmt.Printf("\t%  d. %s\n", i+1, f)
	}
	fmt.Printf("\nselect via file number: ")
	reader := bufio.NewReader(os.Stdin)
	numStr, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	num, err := strconv.Atoi(strings.TrimSpace(numStr))
	if err != nil {
		return "", err
	}
	if num < 1 || num > len(files) {
		return "", fmt.Errorf("%d is not a valid filenumber", num)
	}

	f = files[num-1]
	return f, nil
}

func fullWorkflow(osArgs []string) (args []string, err error) {
	input, err := ExpandArgs(osArgs)
	if err == nil || err != NoArgs {
		return args, err
	}

	appname := input[0]
	input, hasNoAutoArgs := PullNoAutoArgs(input)
	input, hasNoArgs := PullNoArgs(input)
	dirfiles, err := GetDirFiles(path.Base(appname))
	args = []string{appname}
	if len(dirfiles.DefArgs) > 0 && !hasNoArgs && !hasNoAutoArgs {
		args = append(args, "--args", dirfiles.DefArgs)
		args = append(args, input[1:]...)
		return ExpandArgs(args)
	} else if len(dirfiles.ArgsFiles) > 0 && !hasNoArgs {
		file, err := selectArgsFile(dirfiles.ArgsFiles)
		if err != nil {
			return args, err
		}
		args = append(args, "--args", file)
		args = append(args, input[1:]...)
		return ExpandArgs(args)
	}

	return args, nil
}

//******************************************************************************

// NoArgs is the error returned by ExpandArgs when no '--args' argument exist.
var NoArgs = errors.New("NoArgs")

// PullNoAutoArgs pulls the argument '--no-auto-args' from args.
// The return value found will be true if args is reduced.
func PullNoAutoArgs(args []string) (reduced []string, found bool) {
	idx, found := stringIndex(args, "--no-auto-args")
	if !found {
		return args, found
	}
	reduced = append(args[:idx], args[idx+1:]...)
	return reduced, found
}

// PullNoArgs pulls the argument '--no-args' from args.
// The return value found will be true if args is reduced.
func PullNoArgs(args []string) (reduced []string, found bool) {
	idx, found := stringIndex(args, "--no-args")
	if !found {
		return args, found
	}
	reduced = append(args[:idx], args[idx+1:]...)
	return reduced, found
}

// Args implements the most user-friendly integration to handle args files for a
// CLI tool.
// The implementation supports default arguments for an application if no --args
// argument is used.
// Also scans the implementation for all args files in a working directory and
// allow the user to select the args file.
func Args() (args []string, err error) {
	return fullWorkflow(os.Args)
}

// ExpandArgs expands args with the value in the specified args file.
// It returns the error NoArgs if no '--args' argument exist in args.
func ExpandArgs(args []string) (expanded []string, err error) {
	idx, ok := stringIndex(args, "--args")
	if !ok {
		return args, NoArgs
	}

	if idx == len(args)-1 {
		return args, fmt.Errorf("missing --args filepath value")
	}

	expanded = make([]string, idx)
	copy(expanded, args[:idx])
	{
		inserted, err := ReadFile(args[idx+1])
		if err != nil {
			return expanded, err
		}
		expanded = append(expanded, inserted...)
	}
	expanded = append(expanded, args[idx+2:]...)

	return expanded, nil
}

// DirFiles represents the args files in a directory.
type DirFiles struct {
	DefArgs   string
	ArgsFiles []string
}

// GetDirFiles returns the args files for appname in the working directory.
func GetDirFiles(appname string) (dirfiles DirFiles, err error) {
	wd, err := os.Getwd()
	if err != nil {
		return dirfiles, err
	}

	return GetDirFilesIn(appname, wd)
}

// GetDirFilesIn returns the args files for appname in the directory dir.
func GetDirFilesIn(appname string, dir string) (dirfiles DirFiles, err error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return dirfiles, err
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if file.Name() == appname+".auto.args" {
			dirfiles.DefArgs = file.Name()
		} else if strings.HasPrefix(file.Name(), appname) &&
			strings.HasSuffix(file.Name(), ".args") {
			dirfiles.ArgsFiles = append(dirfiles.ArgsFiles, file.Name())
		}
	}

	return dirfiles, nil
}

// Read reads all args from r.
// A successful call returns err == nil, not err == io.EOF.
// Because Read is defined to read until EOF, it does not treat end of file as
// an error to be reported.
func Read(r io.Reader) (args []string, err error) {

	reader := bufio.NewReader(r)

	for {
		line, isPrefix, err := reader.ReadLine()
		if err != nil && err != io.EOF {
			return args, err
		}
		if isPrefix {
			return args, fmt.Errorf("to large line")
		}

		lineStr := string(line)
		prefixIs := func(p string) bool { return strings.HasPrefix(lineStr, p) }
		withoutPrefix := func(p string) string {
			return strings.TrimPrefix(lineStr, p)
		}

		if len(strings.TrimSpace(lineStr)) == 0 {
			// ignore empty line
		} else if prefixIs("#") {
			// ignore comment
		} else if prefixIs("$ ") {
			values := parseShellLine(withoutPrefix("$ "))
			args = append(args, values...)
		} else if prefixIs("| ") {
			args = appendToLast(args, "", withoutPrefix("| "))
		} else if prefixIs("|= ") {
			args = appendToLast(args, "=", withoutPrefix("|= "))
		} else if prefixIs("|s ") {
			args = appendToLast(args, " ", withoutPrefix("|s "))
		} else if prefixIs("|t ") {
			args = appendToLast(args, "\t", withoutPrefix("|t "))
		} else if prefixIs("|n ") {
			args = appendToLast(args, "\n", withoutPrefix("|n "))
		} else {
			args = append(args, lineStr)
		}

		if err != nil {
			return args, nil
		}
	}
}

// ReadFile reads all args from filepath.
// A successful call returns err == nil, not err == io.EOF.
// Because Read is defined to read until EOF, it does not treat end of file as
// an error to be reported.
func ReadFile(filepath string) (args []string, err error) {
	file, err := os.Open(filepath)
	if err != nil {
		return args, err
	}
	defer file.Close()

	return Read(file)
}
