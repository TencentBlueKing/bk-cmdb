// Jsontoml reads JSON and converts to TOML.
//
// Usage:
//   cat file.toml | jsontoml > file.json
//   jsontoml file1.toml > file.json
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/pelletier/go-toml"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "jsontoml can be used in two ways:")
		fmt.Fprintln(os.Stderr, "Writing to STDIN and reading from STDOUT:")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Reading from a file name:")
		fmt.Fprintln(os.Stderr, " tomljson file.toml")
	}
	flag.Parse()
	os.Exit(processMain(flag.Args(), os.Stdin, os.Stdout, os.Stderr))
}

func processMain(files []string, defaultInput io.Reader, output io.Writer, errorOutput io.Writer) int {
	// read from stdin and print to stdout
	inputReader := defaultInput

	if len(files) > 0 {
		file, err := os.Open(files[0])
		if err != nil {
			printError(err, errorOutput)
			return -1
		}
		inputReader = file
		defer file.Close()
	}
	s, err := reader(inputReader)
	if err != nil {
		printError(err, errorOutput)
		return -1
	}
	io.WriteString(output, s)
	return 0
}

func printError(err error, output io.Writer) {
	io.WriteString(output, err.Error()+"\n")
}

func reader(r io.Reader) (string, error) {
	jsonMap := make(map[string]interface{})
	jsonBytes, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(jsonBytes, &jsonMap)
	if err != nil {
		return "", err
	}

	tree, err := toml.TreeFromMap(jsonMap)
	if err != nil {
		return "", err
	}
	return mapToTOML(tree)
}

func mapToTOML(t *toml.Tree) (string, error) {
	tomlBytes, err := t.ToTomlString()
	if err != nil {
		return "", err
	}
	return string(tomlBytes[:]), nil
}
