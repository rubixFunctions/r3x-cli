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
	"encoding/base64"
	"encoding/json"
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
		push, err := cmd.Flags().GetBool("push")
		if err != nil{
			panic(err)
		}
		create(push)

	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().BoolP("push", "p", false, "Push Image")

}

func create(push bool) {
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
	err = tar.Create("/tmp/archieve.tar")
	if  err != nil {
		panic(err)
	}
	err = tar.AddAll(wd, false)
	if err != nil {
		panic(err)
	}
	err = tar.Close()
	if err != nil {
		panic(err)
	}
	dockerBuildContext, err := os.Open("/tmp/archieve.tar")
	defer dockerBuildContext.Close()
	cli, _ := client.NewClientWithOpts(client.FromEnv)
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

	if push {
		authString :=  types.AuthConfig{
			Username: "<<name>>",
			Password: "<<password>>",
			Email: "<<email>>",
		}
		encodedJSON, err := json.Marshal(authString)
		if err != nil {
			panic(err)
		}
		authStr := base64.URLEncoding.EncodeToString(encodedJSON)

		pushOptions := types.ImagePushOptions{
			RegistryAuth: authStr,
		}

		pushResponse, err := cli.ImagePush(context.Background(), "<<name>>/"+funcName, pushOptions)
		if err != nil {
			fmt.Printf("%s", err.Error())
		}
		fmt.Println("Pushing Image has Started")
		termFD, isTErm := term.GetFdInfo(os.Stderr)
		jsonmessage.DisplayJSONMessagesStream(pushResponse, os.Stderr, termFD, isTErm, nil)
	}


}

func getName() string {
	return LoadSchema().Name
}


