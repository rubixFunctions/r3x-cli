// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
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
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:     "init [name]",
	Aliases: []string{"initialize", "initialise", "create"},
	Short:   "Initialize a Function as a Container",
	Long: `Initialize (r3x init) will create a new Function as a container, 
with a license and the appropriate structure needed for a Knative Function.

	- If an absolute path is provided, it will be created.
	
Init will not use an existing directory with contents.`,
	Run: func(cmd *cobra.Command, args []string) {
		wd, err := os.Getwd()
		if err != nil {
			log.Print(err)
		}

		var function *Function
		if len(args) == 0 {
			fmt.Println("Function name needed")
		} else if len(args) == 1 {
			arg := args[0]
			if arg[0] == '.' {
				arg = filepath.Join(wd, arg)
			}
			function = NewFunction(arg)
			initializeFunction(function)
			fmt.Println(`Your Function is ready at` + function.AbsPath())
		}
	},
}

func initializeFunction(function *Function) {
	if !exists(function.AbsPath()) {
		err := os.MkdirAll(function.AbsPath(), os.ModePerm)
		if err != nil {
			fmt.Println(err)
		}
	} else if !isEmpty(function.AbsPath()) {
		fmt.Println("Function can not be bootstrapped in a non empty direcctory: " + function.AbsPath())
	}

	//createJavaScriptFile(function)
	jSTemplate := `const r3x = require('@rubixfunctions/r3x-js-sdk/build/src/r3x')

let schema
r3x.execute(function(){
	let response = {'message' : 'Hello r3x function'}
	return response 
}, schema)`

	dockerTemplate := `FROM node:alpine

WORKDIR /usr/src/app

COPY package*.json ./

RUN npm install --only=production

COPY . .

ENV PORT 8080
EXPOSE $PORT

CMD [ "npm", "start" ]`

	tempPackageTemplate := `{
  "name": "r3x-js-showcase",
  "version": "0.0.1",
  "description": "r3x JS Showcase Application",
  "main": "index.js",
  "scripts": {
    "start": "node r3x-func.js"
  },
  "repository": {
    "type": "git",
    "url": "git+https://github.com/rubixFunctions/r3x-js-showcase.git"
  },
  "keywords": [
    "javascript",
    "knative"
  ],
  "author": "ciaran roche",
  "license": "Apache-2.0",
  "bugs": {
    "url": "https://github.com/rubixFunctions/r3x-js-showcase/issues"
  },
  "homepage": "https://github.com/rubixFunctions/r3x-js-showcase#readme",
  "dependencies": {
    "@rubixfunctions/r3x-js-sdk": "0.0.2"
  }
}
`

	createFile(function, jSTemplate, "r3x-func.js")
	createFile(function, dockerTemplate, "Dockerfile")
	createFile(function, tempPackageTemplate, "package.json")
}

func init() {
	rootCmd.AddCommand(initCmd)
}