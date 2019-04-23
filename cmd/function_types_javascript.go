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
)

// Executes the creation of a JS function
func InitializeJSFunction(function *Function, schema *Schema){
	if !exists(function.AbsPath()) {
		err := os.MkdirAll(function.AbsPath(), os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else if !isEmpty(function.AbsPath()) {
		fmt.Println("Function can not be bootstrapped in a non empty direcctory: " + function.AbsPath())
		return
	}

	createJSDockerfile(function)
	createJSMain(function)
	createSchema(schema, function)
	createJSPackageJSON(function)
	createLicense(function)
	fmt.Println(`Your Function is ready at` + function.AbsPath())
}

// Creates JS specific Dockerfile
func createJSDockerfile(function *Function) {

	dockerTemplate := `FROM node:alpine

WORKDIR /usr/src/app

COPY . .

RUN npm install --only=production

ENV PORT 8080
EXPOSE $PORT

CMD [ "npm", "start" ]`


	createFile(function, dockerTemplate, "Dockerfile")
}

// Creates main JS func file
func createJSMain(function *Function){
	jSTemplate := `const r3x = require('@rubixfunctions/r3x-js-sdk')

/**
 * Simply pass your logic to 'r3x execute'
 */
r3x.execute(function(){
	let response = {'message' : 'Hello r3x function'}
	return response 
})`
	createFile(function, jSTemplate, "r3x-func.js")
}

// Creates JS package.json file
func createJSPackageJSON(function *Function) {
	tempPackageTemplate := `{
		"name": "{{ .name}}",
		"version": "0.0.1",
		"description": "r3x Knative Function",
		"main": "r3x-func.js",
		"scripts": {
		  "start": "node r3x-func.js"
		},
		"keywords": [
		  "javascript",
		  "knative",
		  "kubernetes",
		  "serverless"
		],
		"dependencies": {
		  "@rubixfunctions/r3x-js-sdk": "0.0.14"
		}
	  }
	  `

	data := make(map[string]interface{})
	data["name"] = function.name

	rootCmdScript, err := executeTemplate(tempPackageTemplate, data)
	if err != nil {
		fmt.Println(err)
	}

	err = writeStringToFile(filepath.Join(function.AbsPath(), "package.json"), rootCmdScript)
	if err != nil {
		fmt.Println(err)
	}
}