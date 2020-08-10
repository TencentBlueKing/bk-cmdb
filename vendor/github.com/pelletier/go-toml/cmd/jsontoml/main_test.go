package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"testing"
)

func expectBufferEquality(t *testing.T, name string, buffer *bytes.Buffer, expected string) {
	output := buffer.String()
	if output != expected {
		t.Errorf("incorrect %s: \n%sexpected %s: \n%s", name, output, name, expected)
		t.Log([]rune(output))
		t.Log([]rune(expected))
	}
}

func expectProcessMainResults(t *testing.T, input string, args []string, exitCode int, expectedOutput string, expectedError string) {
	inputReader := strings.NewReader(input)

	outputBuffer := new(bytes.Buffer)
	errorBuffer := new(bytes.Buffer)

	returnCode := processMain(args, inputReader, outputBuffer, errorBuffer)

	expectBufferEquality(t, "output", outputBuffer, expectedOutput)
	expectBufferEquality(t, "error", errorBuffer, expectedError)

	if returnCode != exitCode {
		t.Error("incorrect return code:", returnCode, "expected", exitCode)
	}
}

func TestProcessMainReadFromStdin(t *testing.T) {
	expectedOutput := `
[mytoml]
  a = 42.0
`
	input := `{
  "mytoml": {
    "a": 42
  }
}
`
	expectedError := ``
	expectedExitCode := 0

	expectProcessMainResults(t, input, []string{}, expectedExitCode, expectedOutput, expectedError)
}

func TestProcessMainReadFromFile(t *testing.T) {
	input := `{
  "mytoml": {
    "a": 42
  }
}
`
	tmpfile, err := ioutil.TempFile("", "example.json")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := tmpfile.Write([]byte(input)); err != nil {
		t.Fatal(err)
	}

	defer os.Remove(tmpfile.Name())

	expectedOutput := `
[mytoml]
  a = 42.0
`
	expectedError := ``
	expectedExitCode := 0

	expectProcessMainResults(t, ``, []string{tmpfile.Name()}, expectedExitCode, expectedOutput, expectedError)
}

func TestProcessMainReadFromMissingFile(t *testing.T) {
	var expectedError string
	if runtime.GOOS == "windows" {
		expectedError = `open /this/file/does/not/exist: The system cannot find the path specified.
`
	} else {
		expectedError = `open /this/file/does/not/exist: no such file or directory
`
	}

	expectProcessMainResults(t, ``, []string{"/this/file/does/not/exist"}, -1, ``, expectedError)
}
