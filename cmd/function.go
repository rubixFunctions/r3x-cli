package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// Function type
type Function struct {
	absPath string
	cmdPath string
	srcPath string
	//license License
	name string
}

// NewFunction Returns Function with a specified function name
func NewFunction(functionName string) *Function {
	if functionName == "" {
		fmt.Println("can't create function with no name")
	}

	f := new(Function)
	f.name = functionName

	if f.absPath == "" {
		wd, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
		}
		f.absPath = filepath.Join(wd, functionName)
	}

	fmt.Println(f)
	return f
}

func trimSrcPath(absPath, srcPath string) string {
	relPath, err := filepath.Rel(srcPath, absPath)
	if err != nil {
		fmt.Println(err)
	}
	return relPath
}

func findCmdDir(absPath string) string {
	if !exists(absPath) || isEmpty(absPath) {
		return "cmd"
	}

	if isCmdDir(absPath) {
		return filepath.Base(absPath)
	}

	files, _ := filepath.Glob(filepath.Join(absPath, "c*"))
	for _, file := range files {
		if isCmdDir(file) {
			return filepath.Base(file)
		}
	}

	return "cmd"
}

func isCmdDir(name string) bool {
	name = filepath.Base(name)
	for _, cmdDir := range []string{"cmd", "cmds", "command", "commands"} {
		if name == cmdDir {
			return true
		}
	}
	return false
}

// AbsPath returns Functions absolute path
func (f Function) AbsPath() string {
	return f.absPath
}

// SrcPath returns Functions source path
func (f *Function) SrcPath() string {
	return f.srcPath
}

func filepathHasPrefix(path string, prefix string) bool {
	if len(path) <= len(prefix) {
		return false
	}
	if runtime.GOOS == "windows" {
		return strings.EqualFold(path[0:len(prefix)], prefix)
	}

	return path[0:len(prefix)] == prefix
}
