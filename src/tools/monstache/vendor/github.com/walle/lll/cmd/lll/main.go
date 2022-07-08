package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/alexflint/go-arg"

	"github.com/walle/lll"
)

var args struct {
	MaxLength int      `arg:"-l,env,help:max line length to check for"`
	TabWidth  int      `arg:"-w,env,help:tab width in spaces"`
	GoOnly    bool     `arg:"-g,env,help:only check .go files"`
	SkipTests bool     `arg:"-t,env,help:skip _test.go files"`
	Input     []string `arg:"positional"`
	SkipList  []string `arg:"-s,env,help:list of dirs to skip"`
	Vendor    bool     `arg:"env,help:check files in vendor directory"`
	Files     bool     `arg:"help:read file names from stdin one at each line"`
	Exclude   string   `arg:"-e,env,help:exclude lines that matches this regex"`
}

func main() {
	args.MaxLength = 80
	args.TabWidth = 1
	args.SkipList = []string{".git", "vendor"}
	arg.MustParse(&args)

	var exclude *regexp.Regexp
	if args.Exclude != "" {
		e, err := regexp.Compile(args.Exclude)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error compiling exclude regexp: %s\n", err)
			os.Exit(1)
		}
		exclude = e
	}

	// If we should include the vendor dir, attempt to remove it from the skip list
	if args.Vendor {
		for i, p := range args.SkipList {
			if p == "vendor" {
				args.SkipList = append(args.SkipList[:i], args.SkipList[:i]...)
			}
		}
	}

	// If we should read files from stdin, read each line and process the file
	if args.Files {
		s := bufio.NewScanner(os.Stdin)
		for s.Scan() {
			err := lll.ProcessFile(os.Stdout, s.Text(),
				args.MaxLength, args.TabWidth, exclude)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error processing file: %s\n", err)
			}
		}
		os.Exit(0)
	}

	// Otherwise, walk the inputs recursively
	for _, d := range args.Input {
		err := filepath.Walk(d, func(p string, i os.FileInfo, e error) error {
			if i == nil {
				fmt.Fprintf(os.Stderr, "lll: %s no such file or directory\n", p)
				return nil
			}
			if e != nil {
				fmt.Fprintf(os.Stderr, "lll: %s\n", e)
				return nil
			}
			skip, ret := lll.ShouldSkip(p, i.IsDir(), args.SkipList, args.GoOnly, args.SkipTests)
			if skip {
				return ret
			}

			err := lll.ProcessFile(os.Stdout, p, args.MaxLength, args.TabWidth, exclude)
			return err
		})

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error walking the file system: %s\n", err)
			os.Exit(1)
		}
	}
}
