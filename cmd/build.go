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
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"os"
	"syscall"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/jhoonb/archivex"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build a RubiX Function as a Container",
	Long: `
Build a RubiX Function as a Container Image.
The Image will be pushed to a specified registry
`,
	Run: func(cmd *cobra.Command, args []string) {
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			panic(err)
		}
		if name == ""{
			fmt.Println("UserName or OrgName needed")
			return
		}
		push, err := cmd.Flags().GetBool("push")
		if err != nil{
			panic(err)
		}
		quay, err := cmd.Flags().GetBool("quay")
		if err != nil{
			panic(err)
		}
		create(name, push, quay)
	},
}

// Declare func flags and adds func to root
func init() {
	rootCmd.AddCommand(buildCmd)

	buildCmd.Flags().BoolP("push", "p", false, "Push Image")
	buildCmd.Flags().BoolP("quay", "q", false, "Push to Quay.io")
	buildCmd.Flags().StringP("name", "n", "", "UserName or Org")
}

// Create function, builds tar of directory and execute Build Function
func create(name string, push bool, quay bool) {
	fmt.Println("RubiX is starting to build function image")
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	funcName := getName()
	if funcName == "" {
		fmt.Println("A function needs a name")
		return
	}
	pass := getPass()
	fmt.Println("Creating function tar")
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
	var imageName string
	if quay {
		imageName = "quay.io/" + name + "/" + funcName
	} else {
		imageName = "docker.io/" + name + "/" + funcName
	}
	dockerBuildContext, err := os.Open("/tmp/archieve.tar")
	defer dockerBuildContext.Close()

	cli, _ := client.NewClientWithOpts(client.FromEnv)
	buildImage(imageName, cli, dockerBuildContext)

	if push {
		pushImage(name, pass, imageName, cli)
	}

	genServiceYaml(funcName, imageName)

}

// Passes tar to docker daemon to create function image
func buildImage(imageName string, cli client.ImageAPIClient, dockerBuildContext io.Reader){
	options := types.ImageBuildOptions{
		SuppressOutput: false,
		Remove:         true,
		ForceRemove:    true,
		// hard coded tag, till schema is added to sdk
		Tags:       []string{imageName},
		PullParent: true}
	buildResponse, err := cli.ImageBuild(context.Background(), dockerBuildContext, options)
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
	fmt.Println("Build Image has Started ")
	termFd, isTerm := term.GetFdInfo(os.Stderr)
	jsonmessage.DisplayJSONMessagesStream(buildResponse.Body, os.Stderr, termFd, isTerm, nil)

}

// Pushes Image to Registry
func pushImage(name string, pass string, imageName string, cli client.ImageAPIClient){
	authString :=  types.AuthConfig{
		Username: name,
		Password: pass,
	}
	encodedJSON, err := json.Marshal(authString)
	if err != nil {
		panic(err)
	}
	authStr := base64.URLEncoding.EncodeToString(encodedJSON)

	pushOptions := types.ImagePushOptions{
		RegistryAuth: authStr,
	}

	pushResponse, err := cli.ImagePush(context.Background(), imageName, pushOptions)
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
	fmt.Println("Pushing Image has Started")
	termFD, isTErm := term.GetFdInfo(os.Stderr)
	jsonmessage.DisplayJSONMessagesStream(pushResponse, os.Stderr, termFD, isTErm, nil)
}

// Generates service yaml on successful push to registry
func genServiceYaml(name string, image string){
	schema := LoadSchema()
	switch schema.FuncType {
	case "js":
		createServiceYAML(name, image)
	case "go":
		createServiceYAML(name, image)
	case "py":
		createServiceYAML(name, image)
	case "hs":
		createServiceYAML(name, image)
	default:
		fmt.Println("Error parsing Schema, no service.yaml generated")
	}
}

// Get user password
func getPass() string {
	fmt.Print("Enter Password for Image Hub: ")
	bytePass, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		panic(err)
	}
	return string(bytePass)
}

// Get function name from schema
func getName() string {
	var name = LoadSchema().Name
	if name == ""{
		return ""
	}
	return name
}


