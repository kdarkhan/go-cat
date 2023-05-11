package main

import (
	"bufio"
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSameAsCat(t *testing.T) {

	tests := []struct {
		name      string
		args      string
		isFailure bool
	}{
		{
			name: "Single file",
			args: "testdata/file1.txt",
		},
		{
			name: "Two files",
			args: "testdata/file1.txt testdata/file2.txt",
		},
		{
			name: "Param --show-all",
			args: "--show-all testdata/file1.txt testdata/file2.txt",
		},
		{
			name: "Param --show-all short version",
			args: "-A testdata/file1.txt testdata/file2.txt",
		},
		{
			name: "Param --number-nonblank",
			args: "--number-nonblank testdata/file1.txt testdata/file2.txt",
		},
		{
			name: "Param -e",
			args: "-e testdata/file1.txt testdata/file2.txt",
		},
		{
			name: "Param --show-ends",
			args: "--show-ends testdata/file1.txt testdata/file2.txt",
		},
		{
			name: "Param --show-ends short",
			args: "-E testdata/file1.txt testdata/file2.txt",
		},
		{
			name: "Param --number",
			args: "--number testdata/file1.txt testdata/file2.txt",
		},
		{
			name: "Param --number short",
			args: "-n testdata/file1.txt testdata/file2.txt",
		},
		{
			name: "Param --squeeze-blank",
			args: "--squeeze-blank testdata/file1.txt testdata/file2.txt",
		},
		{
			name: "Param --squeeze-blank short",
			args: "-s testdata/file1.txt testdata/file2.txt",
		},
		{
			name: "Param -t",
			args: "-t testdata/file1.txt testdata/file2.txt",
		},
		{
			name: "Param --show-tabs",
			args: "--show-tabs testdata/file1.txt testdata/file2.txt",
		},
		{
			name: "Param --show-tabs short (-T)",
			args: "-T testdata/file1.txt testdata/file2.txt",
		},
		{
			name: "Param -u (ignored)",
			args: "-u testdata/file1.txt testdata/file2.txt",
		},
		{
			name: "Param --show-nonprintintg",
			args: "--show-nonprinting testdata/file1.txt testdata/file2.txt",
		},
		{
			name:      "With missing file",
			args:      "testdata/file1.txt non-existing-file.txt",
			isFailure: true,
		},
	}

	for _, tc := range tests {
		args := strings.Split(tc.args, " ")

		expectedOut, expectedErr := exec.Command("cat", args...).Output()
		if expectedErr != nil && !tc.isFailure {
			t.Errorf("Failed to run cat %s", expectedErr)
		}

		goRunArgs := []string{"run", "."}
		goRunArgs = append(goRunArgs, args...)
		actualOut, actualErr := exec.Command("go", goRunArgs...).Output()
		if actualErr != nil && !tc.isFailure {
			t.Fatalf("%s: could not execute command", tc.name)
		}

		if tc.isFailure {
			assert.NotNil(t, actualErr)
		}

		assert.Equal(t, string(expectedOut), string(actualOut), "Test <%s> failed", tc.name)
	}
}

func TestSameAsCatWithStdin(t *testing.T) {

	tests := []struct {
		name        string
		file        string
		extraParams string
	}{
		{
			name:        "No param",
			file:        "testdata/file1.txt",
			extraParams: "",
		},
		{
			name:        "With - param",
			file:        "testdata/file1.txt",
			extraParams: "-",
		},
		{
			name:        "With file + stdin param",
			file:        "testdata/file1.txt",
			extraParams: "- testdata/file2.txt",
		},
	}

	for _, tc := range tests {
		args := []string{}
		if tc.extraParams != "" {
			args = strings.Split(tc.extraParams, " ")
		}

		expectedOut := runPipingFileToStdin(t, "cat", tc.file, args)

		goRunArgs := []string{"run", "."}
		goRunArgs = append(goRunArgs, args...)

		actualOut := runPipingFileToStdin(t, "go", tc.file, goRunArgs)

		assert.Equal(t, expectedOut, actualOut, "Test <%s> failed", tc.name)
	}
}

func runPipingFileToStdin(t *testing.T, cmd, inputFile string, args []string) []byte {
	command := exec.Command(cmd, args...) //.Output()

	stdin, err := command.StdinPipe()
	if err != nil {
		t.Fatalf("Could not connect to stdin %s", cmd)
	}
	defer stdin.Close()

	file, err := os.Open(inputFile)
	if err != nil {
		t.Error(err)
	}
	defer file.Close()

	var stdout bytes.Buffer
	command.Stdout = &stdout

	if err = command.Start(); err != nil {
		t.Error(err)
	}
	bufio.NewReader(file).WriteTo(stdin)
	stdin.Close()

	if err = command.Wait(); err != nil {
		t.Error(err)
	}

	return stdout.Bytes()
}
