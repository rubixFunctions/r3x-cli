package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"errors"
)

func init() {
	// Mute commands.
	initCmd.SetOutput(new(bytes.Buffer))
}

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

func TestRubiXJSOutputInitCmd(t *testing.T) {
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

		relPath, err := filepath.Rel(testFunction.AbsPath(), path)
		if err != nil {
			return err
		}
		relPath = filepath.ToSlash(relPath)
		goldenPath := filepath.Join("testData", filepath.Base(path)+".golden")

		fmt.Println(goldenPath)

		files, err := ioutil.ReadDir(testFunction.AbsPath())
		if err != nil {
			t.Fatal(err)
		}

		for _, f := range files {
			gotFiles = append(gotFiles, f.Name())
		}

		switch relPath {
		// Known directories.
		case ".", "cmd":
			return nil
		// Known files.
		case "Dockerfile", "LICENSE", "package.json", "r3x-func.js", "schema.json":
			return compareFiles(path, goldenPath)
		}
		// Unknown file.
		return errors.New("unknown file: " + path)
	})
	if err != nil {
		if err := checkLackFiles(expectedFiles, gotFiles); err != nil {
			t.Fatal(err)
		}
	}

}



func TestRubiXGoInitCmd(t *testing.T) {
	functionName := "testFunction"

	testFunction := NewFunction(functionName)
	defer os.RemoveAll(testFunction.AbsPath())

	os.Args = []string{"r3x", "init", functionName, "--type", "go"}
	if err := rootCmd.Execute(); err != nil {
		t.Fatal("Error by execution:", err)
	}

	expectedFiles := []string{"Dockerfile", "LICENSE", "main.go", "schema.json"}
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

func TestRubiXPyInitCmd(t *testing.T) {
	functionName := "testFunction"

	testFunction := NewFunction(functionName)
	defer os.RemoveAll(testFunction.AbsPath())

	os.Args = []string{"r3x", "init", functionName, "--type", "py"}
	if err := rootCmd.Execute(); err != nil {
		t.Fatal("Error by execution:", err)
	}

	expectedFiles := []string{"Dockerfile", "LICENSE", "r3x-func.py", "schema.json"}
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


// ensureLF converts any \r\n to \n
func ensureLF(content []byte) []byte {
	return bytes.Replace(content, []byte("\r\n"), []byte("\n"), -1)
}

// compareFiles compares the content of files with pathA and pathB.
// If contents are equal, it returns nil.
// If not, it returns which files are not equal
// and diff (if system has diff command) between these files.
func compareFiles(pathA, pathB string) error {
	contentA, err := ioutil.ReadFile(pathA)
	if err != nil {
		return err
	}
	contentB, err := ioutil.ReadFile(pathB)
	if err != nil {
		return err
	}
	if !bytes.Equal(ensureLF(contentA), ensureLF(contentB)) {
		output := new(bytes.Buffer)
		output.WriteString(fmt.Sprintf("%q and %q are not equal!\n\n", pathA, pathB))

		diffPath, err := exec.LookPath("diff")
		if err != nil {
			// Don't execute diff if it can't be found.
			return nil
		}
		diffCmd := exec.Command(diffPath, "-u", pathA, pathB)
		diffCmd.Stdout = output
		diffCmd.Stderr = output

		output.WriteString("$ diff -u " + pathA + " " + pathB + "\n")
		if err := diffCmd.Run(); err != nil {
			output.WriteString("\n" + err.Error())
		}
		return errors.New(output.String())
	}
	return nil
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
