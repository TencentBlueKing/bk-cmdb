// Package lll provides validation functions regarding line length
package lll

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode/utf8"
)

// ShouldSkip checks the input and determines if the path should be skipped.
// Use the SkipList to quickly skip paths.
// All directories are skipped, only files are processed.
// If GoOnly is supplied check that the file is a go file.
// Otherwise check so the file is a "text file".
func ShouldSkip(path string, isDir bool, skipList []string,
	goOnly bool, skipTests bool) (bool, error) {

	name := filepath.Base(path)
	for _, d := range skipList {
		if name == d {
			if isDir {
				return true, filepath.SkipDir
			}
			return true, nil
		}
	}
	if isDir {
		return true, nil
	}

	if skipTests && strings.HasSuffix(path, "_test.go") {
		return true, nil
	}
	isGo := strings.HasSuffix(path, ".go")
	if goOnly && !isGo {
		return true, nil
	}
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return true, err
	}
	if isGo {
		return isGenerated(b), nil
	}
	m := http.DetectContentType(b)
	if !strings.Contains(m, "text/") {
		return true, nil
	}
	return false, nil
}

// ProcessFile checks all lines in the file and writes an error if the line
// length is greater than MaxLength.
func ProcessFile(w io.Writer, path string, maxLength, tabWidth int,
	exclude *regexp.Regexp) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() {
		err := f.Close()
		if err != nil {
			fmt.Printf("Error closing file: %s\n", err)
		}
	}()

	return Process(f, w, path, maxLength, tabWidth, exclude)
}

// Process checks all lines in the reader and writes an error if the line length
// is greater than MaxLength.
func Process(r io.Reader, w io.Writer, path string, maxLength, tabWidth int,
	exclude *regexp.Regexp) error {
	spaces := strings.Repeat(" ", tabWidth)
	l := 0
	s := bufio.NewScanner(r)
	for s.Scan() {
		l++
		t := s.Text()
		t = strings.Replace(t, "\t", spaces, -1)
		c := utf8.RuneCountInString(t)
		if c > maxLength {
			if exclude != nil {
				if exclude.MatchString(t) {
					continue
				}
			}
			fmt.Fprintf(w, "%s:%d: line is %d characters\n", path, l, c)
		}
	}

	return s.Err()
}
