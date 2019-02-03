# RubiX Command Line Interface
[![CircleCI](https://circleci.com/gh/rubixFunctions/r3x-cli.svg?style=svg)](https://circleci.com/gh/rubixFunctions/r3x-cli)
[![License](https://img.shields.io/badge/-Apache%202.0-blue.svg)](https://opensource.org/s/Apache-2.0)

# Usage
## Prerequisite
The CLI requires the following :
- GoLang
- Git

## Build Steps
Clone project
```
$ mkdir $GOPATH/src/github.com/rubixFunctions
$ cd $GOPATH/src/github.com/rubixFunctions
$ git clone git@github.com:rubixFunctions/r3x-cli.git
$ cd r3x-cli
$ go build
```
After building the project, add the binary `r3x-cli` to $PATH

## Using CLI
To get help :
```
$ r3x help
```
### Bootstrap a Function as a Container
```
$ r3x init hello-function --type js 
```
**NOTE** Function will be initialized with an Apache 2.0 License as standard, if you require a different License use the `license` tag :
```
$ r3x init hello-function --type js --license MIT
```
Alternativly you can initalize a function with no License by :
```
$ r3x init hello-function --type js --license none
```
### Create a Function as a Container Image
The create function image locally:
```
$ cd <<function dir>>
$ r3x create -n <<Repo Name or Org>>
```
To create and push an image to Docker Hub:
```
$ r3x create -p -n <<Repo Name or Org>>

```
To create and push an image to Quay.io
```
$ r3x create -p -q -n <<Repo Name or Org>>
```

## License
This project is licensed under the Apache 2.0 License - see the [LICENSE](LICENSE) file for details

