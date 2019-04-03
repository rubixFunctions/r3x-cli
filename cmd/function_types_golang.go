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
)

// Executes the creation of the Golang Function
func InitializeGoFunction(function *Function, schema *Schema){
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

	createGoDockerfile(function)
	createGoMain(function)
	createSchema(schema, function)
	createLicense(function)
	fmt.Println(`Your Function is ready at` + function.AbsPath())
}

// Creates go specific Docker file
func createGoDockerfile(function *Function) {
	dockerTemplate := `FROM golang:alpine as builder
RUN apk update && apk add git
RUN adduser -D -g '' appuser
RUN go get github.com/rubixFunctions/r3x-golang-sdk
ENV SOURCES /go/src/github.com/rubixFunctions/r3x-showcase-apps/r3x-golang-showcase
COPY . ${SOURCES}
WORKDIR ${SOURCES}
RUN cd ${SOURCES} && CGO_ENABLED=0 go build
FROM scratch
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /go/src/github.com/rubixFunctions/r3x-showcase-apps/r3x-golang-showcase/r3x-golang-showcase /r3x-golang-showcase
USER appuser
ENV PORT 8080
ENTRYPOINT [ "/r3x-golang-showcase" ]`

	createFile(function, dockerTemplate, "Dockerfile")
}

// Creates go main func file
func createGoMain(function *Function){
	goTemplate := `package main

import (
	"github.com/rubixFunctions/r3x-golang-sdk"
)


func main() {
	r3x.Execute(r3xFunc)
}

func r3xFunc(input map[string]interface{}) []byte {
	// Input from request body can be accessed as a map eg -- name := input["name"]
	// Return function response as byte array
	var response = "{'message': 'hello r3x'}"
	return []byte(response)
}`

	createFile(function, goTemplate, "main.go")

}

