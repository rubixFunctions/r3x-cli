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
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)


// Checks if Path exists
func exists(path string) bool {
	if path == "" {
		return false
	}
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if !os.IsNotExist(err) {
		fmt.Println(err)
	}
	return false
}

// Checks if Path is empty
func isEmpty(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		fmt.Println(err)
	}

	if !fi.IsDir() {
		return fi.Size() == 0
	}

	f, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	names, err := f.Readdirnames(-1)
	if err != nil && err != io.EOF {
		fmt.Println(err)
	}

	for _, name := range names {
		if len(name) > 0 && name[0] != '.' {
			return false
		}
	}
	return true
}

// Executes a file template
func executeTemplate(tmplStr string, data interface{}) (string, error) {
	tmpl, err := template.New("").Funcs(template.FuncMap{"comment": commentifyString}).Parse(tmplStr)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, data)
	return buf.String(), err
}

// Writes contents to a file
func writeToFile(path string, r io.Reader) error {
	if exists(path) {
		return fmt.Errorf("%v already exists", path)
	}

	dir := filepath.Dir(path)
	if dir != "" {
		if err := os.MkdirAll(dir, 0777); err != nil {
			return err
		}
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, r)
	return err
}

// Executes write to file function
func writeStringToFile(path string, s string) error {
	return writeToFile(path, strings.NewReader(s))
}

// Formats new lines in file
func commentifyString(in string) string {
	var newlines []string
	lines := strings.Split(in, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "//") {
			newlines = append(newlines, line)
		} else {
			if line == "" {
				newlines = append(newlines, "//")
			} else {
				newlines = append(newlines, "// "+line)
			}
		}
	}
	return strings.Join(newlines, "\n")
}

// Creates a file
func createFile(function *Function, template, file string) {
	data := make(map[string]interface{})
	rootCmdScript, err := executeTemplate(template, data)
	if err != nil {
		fmt.Println(err)
	}
	err = writeStringToFile(filepath.Join(function.AbsPath(), file), rootCmdScript)
	if err != nil {
		fmt.Println(err)
	}
}

// Creates a Schema
func createSchema(schema *Schema, function *Function){
	data := make(map[string]interface{})
	data["name"] = schema.Name
	data["funcType"] = schema.FuncType
	data["response"] = schema.Response
	var schemaJson = `{
"name" : "{{.name}}",
"funcType" : "{{.funcType}}",
"response" : "{{.response}}"
}`

	rootCmdScript, err := executeTemplate(schemaJson, data)
	if err != nil {
		fmt.Println(err)
	}

	err = writeStringToFile(filepath.Join(function.AbsPath(), "schema.json"), rootCmdScript)
	if err != nil {
		fmt.Println(err)
	}
}

// Creates r3x service YAML
func createServiceYAML(name string, image string){
	wd, err := os.Getwd()
	if err != nil {
		log.Print(err)
		return
	}
	data := make(map[string]interface{})
	data["name"] = name
	data["image"] = image
	var serviceYaml = `apiVersion: serving.knative.dev/v1alpha1
kind: Service
metadata:
  name: {{.name}}
  namespace: default
spec:
  runLatest:
    configuration:
      revisionTemplate:
        spec:
          container:
            image: {{.image}}`

	rootCmdScript, err := executeTemplate(serviceYaml, data)
	if err != nil {
		fmt.Println(err)
	}

	err = writeStringToFile(filepath.Join(wd, "service.yaml"), rootCmdScript)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Push completed and service.yaml generated")
}

var logo =`
		______      _     ___   __
		| ___ \    | |   (_) \ / /
		| |_/ /   _| |__  _ \ V / 
		|    / | | | '_ \| |/   \ 
		| |\ \ |_| | |_) | / /^\ \
		\_| \_\__,_|_.__/|_\/   \/

`

func initString(t string) {
	fmt.Println(fmt.Sprintf("Building %v Scaffolding", t))
}
