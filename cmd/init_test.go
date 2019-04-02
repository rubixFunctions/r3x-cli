package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestRubiXJSInitCmd(t *testing.T) {
	functionName := "testFunction"

	testFunction := NewFunction(functionName)
	defer os.RemoveAll(testFunction.AbsPath())

	os.Args = []string{"r3x", "init", functionName, "--type", "js"}
	if err := rootCmd.Execute(); err != nil {
		t.Fatal("Error by execution:", err)
	}

	expectedFiles := []string{"Dockerfile", "LICENSE", "package.json", "r3x-func.js", "schema.json"}
	gotFiles := []string{}

	err := filepath.Walk(testFunction.AbsPath(), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		files, err := ioutil.ReadDir(testFunction.AbsPath())
		if err != nil {
			t.Fatal(err)
		}

		for _, f := range files {
			gotFiles = append(gotFiles, f.Name())
		}

		return checkLackFiles(expectedFiles, gotFiles)
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestRubiXNoLicenseInitCmd(t *testing.T){
	functionName := "testFunction"

	testFunction := NewFunction(functionName)
	defer os.RemoveAll(testFunction.AbsPath())

	os.Args = []string{"r3x", "init", functionName, "--type", "js", "--license", "none"}
	if err := rootCmd.Execute(); err != nil {
		t.Fatal("Error by execution:", err)
	}

	expectedFiles := []string{"Dockerfile", "package.json", "r3x-func.js", "schema.json"}
	gotFiles := []string{}

	err := filepath.Walk(testFunction.AbsPath(), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		files, err := ioutil.ReadDir(testFunction.AbsPath())
		if err != nil {
			t.Fatal(err)
		}

		for _, f := range files {
			gotFiles = append(gotFiles, f.Name())
		}

		return checkLackFiles(expectedFiles, gotFiles)
	})
	if err != nil {
		t.Fatal(err)
	}
}

// checkLackFiles checks if all elements of expected are in got.
func checkLackFiles(expected, got []string) error {
	lacks := make([]string, 0, len(expected))
	for _, ev := range expected {
		if !stringInStringSlice(ev, got) {
			lacks = append(lacks, ev)
		}
	}
	if len(lacks) > 0 {
		return fmt.Errorf("Lack %v file(s): %v", len(lacks), lacks)
	}
	return nil
}

// stringInStringSlice checks if s is an element of slice.
func stringInStringSlice(s string, slice []string) bool {
	for _, v := range slice {
		if s == v {
			return true
		}
	}
	return false
}
