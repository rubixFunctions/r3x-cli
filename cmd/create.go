// Copyright © 2018 RubixFunctions.
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
	"context"
	"fmt"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/term"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/jhoonb/archivex"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Build a RubiX Function as a Container",
	Long: `
Build a RubiX Function as a Container Image.
The Image will be pushed to a specified registery
`,
	Run: func(cmd *cobra.Command, args []string) {
		create()
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

}

func create() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	funcName := getName()
	if funcName == "" {
		panic("A function needs a name")
		return
	}
	tar := new(archivex.TarFile)
	_ = tar.Create(wd + "/archieve.tar")
	_ = tar.AddAll(wd, false)
	_ = tar.Close()
	dockerBuildContext, err := os.Open(wd + "/archieve.tar")
	defer dockerBuildContext.Close()
	defaultHeaders := map[string]string{"User-Agent": "ego-v-0.0.1"}
	cli, _ := client.NewClient("unix:///var/run/docker.sock", "v1.24", nil, defaultHeaders)
	options := types.ImageBuildOptions{
		SuppressOutput: false,
		Remove:         true,
		ForceRemove:    true,
		// hard coded tag, till schema is added to sdk
		Tags:       []string{funcName},
		PullParent: true}
	buildResponse, err := cli.ImageBuild(context.Background(), dockerBuildContext, options)
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
	fmt.Println("Build Image has Started ")
	termFd, isTerm := term.GetFdInfo(os.Stderr)
	jsonmessage.DisplayJSONMessagesStream(buildResponse.Body, os.Stderr, termFd, isTerm, nil)
}
func getName() string {
	return LoadSchema().Name
}
