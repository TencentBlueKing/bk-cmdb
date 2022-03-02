package lll_test

import (
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/walle/lll"
)

func TestShouldSkipDirs(t *testing.T) {
	skip, err := lll.ShouldSkip(".git", true, []string{".git"}, false, false)
	if skip == false || err != filepath.SkipDir {
		t.Errorf("Expected %t, %s got. %t, %s", true, filepath.SkipDir, skip, err)
	}

	skip, err = lll.ShouldSkip("dir", true, []string{".git"}, false, false)
	if skip == false || err != nil {
		t.Errorf("Expected %t, %v got. %t, %s", true, nil, skip, err)
	}
}

func TestShouldSkipFiles(t *testing.T) {
	t.Run("regular files", func(t *testing.T) {
		binaryFilePath, _ := os.Executable()
		tests := []struct {
			path      string
			goOnly    bool
			skipTests bool

			shouldSkip bool
			err        error
		}{
			{path: "lll.go", goOnly: false, skipTests: false, shouldSkip: false},
			{path: "lll.go", goOnly: true, skipTests: false, shouldSkip: false},
			{path: "lll.go", goOnly: false, skipTests: true, shouldSkip: false},
			{path: "lll.go", goOnly: true, skipTests: true, shouldSkip: false},
			{path: "README.md", goOnly: false, skipTests: false, shouldSkip: false},
			{path: "README.md", goOnly: true, skipTests: false, shouldSkip: true},
			{path: "README.md", goOnly: false, skipTests: true, shouldSkip: false},
			{path: "README.md", goOnly: true, skipTests: true, shouldSkip: true},
			{path: "lll_test.go", goOnly: false, skipTests: false, shouldSkip: false},
			{path: "lll_test.go", goOnly: true, skipTests: false, shouldSkip: false},
			{path: "lll_test.go", goOnly: false, skipTests: true, shouldSkip: true},
			{path: "lll_test.go", goOnly: true, skipTests: true, shouldSkip: true},
			{path: binaryFilePath, goOnly: false, skipTests: false, shouldSkip: true},
			{path: binaryFilePath, goOnly: true, skipTests: false, shouldSkip: true},
			{path: binaryFilePath, goOnly: false, skipTests: true, shouldSkip: true},
			{path: binaryFilePath, goOnly: true, skipTests: true, shouldSkip: true},
			{path: "generated_test.go", goOnly: false, skipTests: false, shouldSkip: true},
			{path: "generated_test.go", goOnly: true, skipTests: false, shouldSkip: true},
			{path: "generated_test.go", goOnly: false, skipTests: true, shouldSkip: true},
			{path: "generated_test.go", goOnly: true, skipTests: true, shouldSkip: true},
		}

		for i, tc := range tests {
			skip, err := lll.ShouldSkip(tc.path, false, []string{".git"}, tc.goOnly, tc.skipTests)
			if skip != tc.shouldSkip || err != tc.err {
				t.Errorf("%d) Expected %t, %v got %t, %s", i+1, tc.shouldSkip, tc.err, skip, err)
			}
		}
	})
	t.Run("file in skiplist", func(t *testing.T) {
		skip, err := lll.ShouldSkip("file", false, []string{"file"}, false, false)
		if skip != true || err != nil {
			t.Errorf("Expected %t, %v got. %t, %s", true, nil, skip, err)
		}
	})
	t.Run("error on file not found", func(t *testing.T) {
		skip, err := lll.ShouldSkip("file", false, []string{".git"}, false, false)
		if skip != true || err == nil {
			t.Errorf("Expected %t, %v got. %t, %s", true, nil, skip, err)
		}
	})
}

func TestProcess(t *testing.T) {
	lines := "one\ntwo\ntree"
	b := bytes.NewBufferString("")
	err := lll.Process(bytes.NewBufferString(lines), b, "file", 80, 1, nil)
	if err != nil {
		t.Errorf("Expected %v, got %s", nil, err)
	}

	expected := "file:3: line is 4 characters\n"
	_ = lll.Process(bytes.NewBufferString(lines), b, "file", 3, 1, nil)
	if b.String() != expected {
		t.Errorf("Expected %s, got %s", expected, b.String())
	}
}

func TestProcessFile(t *testing.T) {
	b := bytes.NewBufferString("")
	err := lll.ProcessFile(b, "lll_test.go", 80, 1, nil)
	if err != nil {
		t.Errorf("Expected %v, got %s", nil, err)
	}
}

func TestProcessUnicode(t *testing.T) {
	lines := "日本語\n"
	b := bytes.NewBufferString("")
	expected := "file:1: line is 3 characters\n"
	_ = lll.Process(bytes.NewBufferString(lines), b, "file", 2, 1, nil)
	if b.String() != expected {
		t.Errorf("Expected %s, got %s", expected, b.String())
	}
}

func TestProcessWithTabwidth4(t *testing.T) {
	lines := "\t\t\terr := lll.ProcessFile(os.Stdout, s.Text(), " +
		"args.MaxLength, args.TabWidth, exclude)"
	b := bytes.NewBufferString("")
	expected := "file:1: line is 95 characters\n"
	_ = lll.Process(bytes.NewBufferString(lines), b, "file", 80, 4, nil)
	if b.String() != expected {
		t.Errorf("Expected %s, got %s", expected, b.String())
	}
}

func TestProcessWithTabwidth8(t *testing.T) {
	lines := "\t\t\terr := lll.ProcessFile(os.Stdout, s.Text(), " +
		"args.MaxLength, args.TabWidth, exclude)"
	b := bytes.NewBufferString("")
	expected := "file:1: line is 107 characters\n"
	_ = lll.Process(bytes.NewBufferString(lines), b, "file", 80, 8, nil)
	if b.String() != expected {
		t.Errorf("Expected %s, got %s", expected, b.String())
	}
}

func TestProcessExclude(t *testing.T) {
	lines := "one\ntwo\nTODO: fix\ntree\nFIXME: do this"
	b := bytes.NewBufferString("")
	exclude := regexp.MustCompile("TODO|FIXME")
	expected := "file:4: line is 4 characters\n"
	_ = lll.Process(bytes.NewBufferString(lines), b, "file", 3, 1, exclude)
	if b.String() != expected {
		t.Errorf("Expected %s, got %s", expected, b.String())
	}
}

func BenchmarkProcessExclude(b *testing.B) {
	lines := "one\ntwo\ntree\nTODO: fix\nFIXME: do this"
	exclude := regexp.MustCompile("TODO|FIXME")
	expected := "file:3: line is 4 characters\n"

	for i := 0; i < b.N; i++ {
		buf := bytes.NewBufferString("")
		_ = lll.Process(bytes.NewBufferString(lines), buf, "file", 3, 1, exclude)
		if buf.String() != expected {
			b.Errorf("Expected %s, got %s", expected, buf.String())
		}
	}
}
