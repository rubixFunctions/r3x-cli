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
	"strings"
)

const (
	LicenseMIT      = "MIT"
	LicenseApache20 = "Apache-2.0"
)

var Licenses = make(map[string]License)

var KnownLicenses = []string{
	LicenseMIT,
	LicenseApache20,
}

type License struct {
	Name            string
	PossibleMatches []string
	Text            string
}

func init() {
	initApache2()
	initMit()
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
	// save knownlicenses to map
	l := make(map[string]bool)
	for i := 0; i < len(KnownLicenses); i++ {
		l[KnownLicenses[i]] = true
	}
	// check is license exists
	if _, ok := l[name]; ok {
		return true
	}
	return false
}

// func matchLicense(name string) string {

// }
