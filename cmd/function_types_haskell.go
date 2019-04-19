package cmd

import (
	"fmt"
	"os"
	"path/filepath"
)

func InitializeHaskellFunction(function *Function, schema *Schema){
	if !exists(function.AbsPath()) {
		err := os.MkdirAll(function.AbsPath(), os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else if !isEmpty(function.AbsPath()) {
		fmt.Println("Function can not be bootstrapped in a non empty directory: " + function.AbsPath())
		return
	}

	createHsDockerfile(function)
	createHsMain(function)
	createHsPackageYaml(function)
	createHsStackFile(function)
	createLicense(function)
	createSchema(schema, function)
	fmt.Println(`Your Function is ready at` + function.AbsPath())

}

func createHsDockerfile(function *Function){
	dockerTemplate := `FROM haskell:8.2.2 as builder

# Add dependiencies needed to builder
RUN apt-get update && apt-get install --yes \
    xz-utils \ 
    build-essential \ 
    libtool \
    libpcre3-dev \
    libpcre3 \
    make

# Copy local code to the container image.
WORKDIR /app
COPY . .

# Build and test our code, then build the “helloworld-haskell-exe” executable.
RUN stack setup
RUN stack build --copy-bins --ghc-options '-static -optl-static -optl-pthread -fPIC'

# Use a Docker multi-stage build to create a lean production image.
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM fpco/haskell-scratch

# Copy the "helloworld-haskell-exe" executable from the builder stage to the production image.
WORKDIR /root/
COPY --from=builder /root/.local/bin/helloworld-haskell-exe .

# Service must listen to $PORT environment variable.
# This default value facilitates local development.
ENV PORT 8080

# Run the web service on container startup.
CMD ["./helloworld-haskell-exe"]`

	createFile(function, dockerTemplate, "Dockerfile")

}

func createHsPackageYaml(function *Function){
	packageTemplate := `name:                helloworld-haskell
version:             0.1.0.0
dependencies:
- base >= 4.7 && < 5
- r3x-haskell-sdk
- wai
- text
- aeson

executables:
  helloworld-haskell-exe:
    main:                Main.hs
    source-dirs:         app
    ghc-options:
    - -threaded
    - -rtsopts
    - -with-rtsopts=-N`

	createFile(function, packageTemplate, "package.yaml")

}

func createHsStackFile(function *Function){
	stackTemplate := `flags: {}
packages:
- .
extra-deps: 
  - r3x-haskell-sdk-0.1.0.0@sha256:5411ed12947b6cc623bbf25997a2abaf26686e5bd97c18dce6a486ab07187df8
resolver: lts-10.7`

	createFile(function, stackTemplate, "stack.yaml")

}

func createHsMain(function *Function){
	mainTemplate := `{-# language OverloadedStrings #-}
{-# language DeriveAnyClass #-}
{-# language DeriveGeneric #-}
module Main where

import Rubix
import Data.Aeson (ToJSON, FromJSON)
import GHC.Generics (Generic)
import qualified Data.Text as DT
import qualified Network.Wai as NW

data Message = Message 
  {
    message :: DT.Text
  } deriving (Generic, ToJSON, FromJSON)

rubixMessage :: Message
rubixMessage = Message{message="Hello RubiX!!!"}

-- Define a response handler
rubixHandler :: Handler NW.Response
rubixHandler = do
  let res = toResponse $ Json rubixMessage
  return res

-- Start the server
main :: IO ()
main = execute rubixHandler
`
	createMainFile(function, mainTemplate, "Main.hs")

}

// Creates a file
func createMainFile(function *Function, template, file string) {
	data := make(map[string]interface{})
	rootCmdScript, err := executeTemplate(template, data)
	if err != nil {
		fmt.Println(err)
	}
	var path = filepath.Join(function.AbsPath(), "/app")
	err = writeStringToFile(filepath.Join(path, file), rootCmdScript)
	if err != nil {
		fmt.Println(err)
	}
}