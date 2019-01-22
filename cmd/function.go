// Copyright © 2019 RubiXFunctions
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	license License
	name    string
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

	return f
}

func trimSrcPath(absPath, srcPath string) string {
	relPath, err := filepath.Rel(srcPath, absPath)
	if err != nil {
		fmt.Println(err)
	}
	return relPath
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
