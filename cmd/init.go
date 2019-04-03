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
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:     "init [function name]",
	Aliases: []string{"initialize", "initialise", "create"},
	Short:   "Initialize a Function as a Container",
	Long: `Initialize (r3x init) will create a new Function as a container, 
with a license and the appropriate structure needed for a Knative Function.

	- If an absolute path is provided, it will be created.
	
Init will not use an existing directory with contents.`,
	Run: func(cmd *cobra.Command, args []string) {
		// get flag values
		license := cmd.Flag("license").Value.String()
		name := cmd.Flag("type").Value.String()

		if findLicense(license) == true || len(license) == 0 {
			wd, err := os.Getwd()
			if err != nil {
				log.Print(err)
			}
			// Switch on different function type flag
			switch name {
			case "js":
				var function *Function
				if len(args) == 0 {
					fmt.Println("Function name needed")
				} else if len(args) == 1 {
					arg := args[0]
					if arg[0] == '.' {
						arg = filepath.Join(wd, arg)
					}
					function = NewFunction(arg)
					function.license.Name = license
					var schema *Schema
					schema = NewSchema("r3x-"+arg, "js", "json")
					InitializeJSFunction(function, schema)
				}
			case "go":
				var function *Function
				if len(args) == 0 {
					fmt.Println("Function name needed")
				} else if len(args) == 1 {
					arg := args[0]
					if arg[0] == '.' {
						arg = filepath.Join(wd, arg)
					}
					function = NewFunction(arg)
					function.license.Name = license
					var schema *Schema
					schema = NewSchema("r3x-"+arg, "go", "json")
					InitializeGoFunction(function, schema)
				}
			case "py":
				var function *Function
				if len(args) == 0 {
					fmt.Println("Function name needed")
				} else if len(args) == 1 {
					arg := args[0]
					if arg[0] == '.' {
						arg = filepath.Join(wd, arg)
					}
					function = NewFunction(arg)
					function.license.Name = license
					var schema *Schema
					schema = NewSchema("r3x-"+arg, "py", "json")
					InitializePyFunction(function, schema)
				}
			default:
				fmt.Println(warningTypeMessage)
			}
		} else {
			fmt.Println(`License choice not recognized.
				
Please insure license choice matches the following:
		`, strings.Join(getPossibleMatches(), ", "))
		}

	},
}

// Warning message triggered on no type flag passed
var warningTypeMessage = `Function type required, use '-t' flag
	Supported paradigms :
		- JavaScript : '-t js'
		- GoLang : '-t go'
		- Python : '-t py'`

// Init Init Function
func init() {
	rootCmd.AddCommand(initCmd)
}
