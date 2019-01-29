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

import "fmt"

type Schema struct {
	name     string
	funcType string
	response string
}

// NewSchema Returns Schema with a specified function name
func NewSchema(functionName string, funcType string, response string) *Schema {
	if functionName == "" {
		fmt.Println("can't create function with no name")
	}
	if funcType == "" {
		fmt.Println("can't create function with no type")
	}
	s := new(Schema)

	if response == "" {
		s.response = "JSON"
	}

	s.response = response
	s.name = functionName
	s.funcType = funcType

	return s
}
