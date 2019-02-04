// Copyright © 2018 RubixFunctions.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"path/filepath"
	"strings"
)

var Licenses = make(map[string]License)

type License struct {
	Name            string
	PossibleMatches []string
	Text            string
}

func init() {
	initApache2()
	initMit()
	initNone()
	initLgpl()
	initAgpl()
	initBsdClause2()
	initBsdClause3()
	initGpl2()
	initGpl3()
	initCDDL()
	initEclipse()
	initISC()
}

func getLicense(name string) License {
	var key = getKey(name)

	if key != "" {
		return Licenses[key]
	}

	return Licenses["apache"]
}

func getKey(name string) string {
	for key, l := range Licenses {
		for _, match := range l.PossibleMatches {
			if strings.EqualFold(name, match) {
				return key
			}
		}
	}

	return ""
}

func findLicense(name string) bool {
	for _, l := range Licenses {
		for _, match := range l.PossibleMatches {
			if strings.EqualFold(name, match) {
				return true
			}
		}
	}
	return false
}

func getPossibleMatches() []string {
	var PossibleMatchesList []string
	for _, l := range Licenses {
		for _, m := range l.PossibleMatches {
			PossibleMatchesList = append(PossibleMatchesList, m)
		}
	}
	return PossibleMatchesList
}

func createLicense(function *Function) {
	var name = function.license.Name

	if name == "" {
		name = "Apache-2.0"
	}

	var lic = getLicense(name)

	if lic.Text != "" {
		data := make(map[string]interface{})
		rootCmdScript, err := executeTemplate(lic.Text, data)
		if err != nil {
			fmt.Println(err)
		}

		err = writeStringToFile(filepath.Join(function.AbsPath(), "LICENSE"), rootCmdScript)
		if err != nil {
			fmt.Println(err)
		}
	}
}