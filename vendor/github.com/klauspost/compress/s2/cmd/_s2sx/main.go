package main

import (
	"bufio"
	"bytes"
	"embed"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/klauspost/compress/s2"
	"github.com/klauspost/compress/s2/cmd/internal/readahead"
)

const (
	opUnpack = iota + 1
	opUnTar
)

var (
	goos   = flag.String("os", runtime.GOOS, "Destination operating system")
	goarch = flag.String("arch", runtime.GOARCH, "Destination architecture")
	cpu    = flag.Int("cpu", runtime.GOMAXPROCS(0), "Compress using this amount of threads")
	max    = flag.String("max", "1G", "Maximum executable size. Rest will be written to another file.")
	safe   = flag.Bool("safe", false, "Do not overwrite output files")
	stdout = flag.Bool("c", false, "Write all output to stdout. Multiple input files will be concatenated")
	remove = flag.Bool("rm", false, "Delete source file(s) after successful compression")
	quiet  = flag.Bool("q", false, "Don't write any output to terminal, except errors")
	untar  = flag.Bool("untar", false, "Untar on destination")
	help   = flag.Bool("help", false, "Display help")

	version = "(dev)"
	date    = "(unknown)"
)

//go:embed sfx-exe
var embeddedFiles embed.FS

func main() {
	flag.Parse()
	args := flag.Args()
	sz, err := toSize(*max)
	exitErr(err)
	if len(args) == 0 || *help {
		_, _ = fmt.Fprintf(os.Stderr, "s2sx v%v, built at %v.\n\n", version, date)
		_, _ = fmt.Fprintf(os.Stderr, "Copyright (c) 2011 The Snappy-Go Authors. All rights reserved.\n"+
			"Copyright (c) 2021 Klaus Post. All rights reserved.\n\n")
		_, _ = fmt.Fprintln(os.Stderr, `Usage: s2sx [options] file1 file2

Compresses all files supplied as input separately.
Use file name - to read from stdin and write to stdout.
If input is already s2 compressed it will just have the executable wrapped.
Use s2c commandline tool for advanced option, eg 'cat file.txt||s2c -|s2sx - >out.s2sx'
Output files are written as 'filename.s2sx' and with '.exe' for windows targets.
If output is big, an additional file with ".more" is written. This must be included as well.
By default output files will be overwritten.

Wildcards are accepted: testdir/*.txt will compress all files in testdir ending with .txt
Directories can be wildcards as well. testdir/*/*.txt will match testdir/subdir/b.txt

Options:`)
		flag.PrintDefaults()
		dir, err := embeddedFiles.ReadDir("sfx-exe")
		exitErr(err)
		_, _ = fmt.Fprintf(os.Stderr, "\nAvailable platforms are:\n\n")
		for _, d := range dir {
			_, _ = fmt.Fprintf(os.Stderr, " * %s\n", strings.TrimSuffix(d.Name(), ".s2"))
		}

		os.Exit(0)
	}

	opts := []s2.WriterOption{s2.WriterBestCompression(), s2.WriterConcurrency(*cpu), s2.WriterBlockSize(4 << 20)}
	wr := s2.NewWriter(nil, opts...)

	wantPlat := *goos + "-" + *goarch
	exec, err := embeddedFiles.ReadFile(path.Join("sfx-exe", wantPlat+".s2"))
	if os.IsNotExist(err) {
		dir, err := embeddedFiles.ReadDir("sfx-exe")
		exitErr(err)
		_, _ = fmt.Fprintf(os.Stderr, "os-arch %v not available. Available sfx platforms are:\n\n", wantPlat)
		for _, d := range dir {
			_, _ = fmt.Fprintf(os.Stderr, "* %s\n", strings.TrimSuffix(d.Name(), ".s2"))
		}
		_, _ = fmt.Fprintf(os.Stderr, "\nUse -os and -arch to specify the destination platform.")
		os.Exit(1)
	}
	exec, err = ioutil.ReadAll(s2.NewReader(bytes.NewBuffer(exec)))
	exitErr(err)

	written := int64(0)
	if int64(len(exec))+1 >= sz {
		exitErr(fmt.Errorf("max size less than unpacker. Max size must be at least %d bytes", len(exec)+1))
	}
	mode := byte(opUnpack)
	if *untar {
		mode = opUnTar
	}

	// No args, use stdin/stdout
	stdIn := len(args) == 1 && args[0] == "-"
	if *stdout || stdIn {
		// Write exec once to stdout
		_, err = os.Stdout.Write(exec)
		exitErr(err)
		_, err = os.Stdout.Write([]byte{mode})
		exitErr(err)
		written += int64(len(exec) + 1)
	}

	if stdIn {
		// Catch interrupt, so we don't exit at once.
		// os.Stdin will return EOF, so we should be able to get everything.
		signal.Notify(make(chan os.Signal, 1), os.Interrupt)
		isCompressed, rd := isS2Input(os.Stdin)
		if isCompressed {
			_, err := io.Copy(os.Stdout, rd)
			exitErr(err)
			return
		}
		wr.Reset(os.Stdout)
		_, err = wr.ReadFrom(rd)
		printErr(err)
		printErr(wr.Close())
		return
	}

	var files []string
	for _, pattern := range args {
		found, err := filepath.Glob(pattern)
		exitErr(err)
		if len(found) == 0 {
			exitErr(fmt.Errorf("unable to find file %v", pattern))
		}
		files = append(files, found...)
	}

	for _, filename := range files {
		func() {
			var closeOnce sync.Once
			dstFilename := fmt.Sprintf("%s%s", strings.TrimPrefix(filename, ".s2"), ".s2sx")
			if *goos == "windows" {
				dstFilename += ".exe"
			}
			// Input file.
			file, err := os.Open(filename)
			exitErr(err)
			defer closeOnce.Do(func() { file.Close() })
			isCompressed, rd := isS2Input(file)
			if !*quiet {
				if !isCompressed {
					fmt.Print("Compressing ", filename, " -> ", dstFilename, " for ", wantPlat)
				} else {
					fmt.Print("Creating sfx archive ", filename, " -> ", dstFilename, " for ", wantPlat)
				}
			}
			src, err := readahead.NewReaderSize(rd, *cpu+1, 1<<20)
			exitErr(err)
			defer src.Close()
			var out io.Writer
			switch {
			case *stdout:
				out = os.Stdout
			default:
				if *safe {
					_, err := os.Stat(dstFilename)
					if !os.IsNotExist(err) {
						exitErr(errors.New("destination file exists"))
					}
				}
				dstFile, err := os.OpenFile(dstFilename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0777)
				exitErr(err)
				defer dstFile.Close()
				sw := &switchWriter{w: dstFile, left: sz, close: nil}
				sw.fn = func() {
					dstFilename := dstFilename + ".more"
					if *safe {
						_, err := os.Stat(dstFilename)
						if !os.IsNotExist(err) {
							exitErr(fmt.Errorf("destination '%s' file exists", dstFilename))
						}
					}
					dstFile, err := os.OpenFile(dstFilename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0777)
					exitErr(err)
					sw.close = dstFile.Close
					sw.w = dstFile
					sw.left = (1 << 63) - 1
				}
				defer sw.Close()
				bw := bufio.NewWriterSize(sw, 4<<20*2)
				defer bw.Flush()
				out = bw
				_, err = out.Write(exec)
				exitErr(err)
				_, err = out.Write([]byte{mode})
			}
			exitErr(err)
			wc := wCounter{out: out}
			start := time.Now()
			var input int64
			if !isCompressed {
				wr.Reset(&wc)
				defer wr.Close()
				input, err = wr.ReadFrom(src)
				exitErr(err)
				err = wr.Close()
				exitErr(err)
			} else {
				input, err = io.Copy(&wc, src)
				exitErr(err)
			}
			if !*quiet {
				elapsed := time.Since(start)
				mbpersec := (float64(input) / (1024 * 1024)) / (float64(elapsed) / (float64(time.Second)))
				pct := float64(wc.n) * 100 / float64(input)
				fmt.Printf(" %d -> %d [%.02f%%]; %.01fMB/s\n", input, wc.n, pct, mbpersec)
			}
			if *remove {
				closeOnce.Do(func() {
					file.Close()
					if !*quiet {
						fmt.Println("Removing", filename)
					}
					err := os.Remove(filename)
					exitErr(err)
				})
			}
		}()
	}

}

func printErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "\nERROR:", err.Error())
	}
}

func exitErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "\nERROR:", err.Error())
		os.Exit(2)
	}
}

type wCounter struct {
	n   int
	out io.Writer
}

func (w *wCounter) Write(p []byte) (n int, err error) {
	n, err = w.out.Write(p)
	w.n += n
	return n, err

}

func isS2Input(rd io.Reader) (bool, io.Reader) {
	var tmp [4]byte
	n, err := io.ReadFull(rd, tmp[:])
	switch err {
	case nil:
		return bytes.Equal(tmp[:], []byte{0xff, 0x06, 0x00, 0x00}), io.MultiReader(bytes.NewBuffer(tmp[:]), rd)
	case io.ErrUnexpectedEOF, io.EOF:
		return false, bytes.NewBuffer(tmp[:n])
	}
	exitErr(err)
	return false, nil
}

// toSize converts a size indication to bytes.
func toSize(size string) (int64, error) {
	size = strings.ToUpper(strings.TrimSpace(size))
	firstLetter := strings.IndexFunc(size, unicode.IsLetter)
	if firstLetter == -1 {
		firstLetter = len(size)
	}

	bytesString, multiple := size[:firstLetter], size[firstLetter:]
	bytes, err := strconv.ParseInt(bytesString, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("unable to parse size: %v", err)
	}

	switch multiple {
	case "G", "GB", "GIB":
		return bytes * 1 << 30, nil
	case "M", "MB", "MIB":
		return bytes * 1 << 20, nil
	case "K", "KB", "KIB":
		return bytes * 1 << 10, nil
	case "B", "":
		return bytes, nil
	default:
		return 0, fmt.Errorf("unknown size suffix: %v", multiple)
	}
}

type switchWriter struct {
	w     io.Writer
	left  int64
	fn    func()
	close func() error
}

func (w *switchWriter) Write(b []byte) (int, error) {
	if int64(len(b)) <= w.left {
		w.left -= int64(len(b))
		return w.w.Write(b)
	}
	n, err := w.w.Write(b[:w.left])
	if err != nil {
		return n, err
	}
	w.fn()
	n2, err := w.Write(b[n:])
	return n + n2, err
}
func (w *switchWriter) Close() error {
	if w.close == nil {
		return nil
	}
	return w.close()
}
