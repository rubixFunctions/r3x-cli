# RubiX Command Line Interface
[![CircleCI](https://circleci.com/gh/rubixFunctions/r3x-cli.svg?style=svg)](https://circleci.com/gh/rubixFunctions/r3x-cli)
[![License](https://img.shields.io/badge/-Apache%202.0-blue.svg)](https://opensource.org/s/Apache-2.0)
[![Coverage Status](https://coveralls.io/repos/github/rubixFunctions/r3x-cli/badge.svg?branch=master)](https://coveralls.io/github/rubixFunctions/r3x-cli?branch=master)

## Documentation
For full information on how to use the SDK and deploy a function to Knative, refer to our [Documentation here.](https://github.com/rubixFunctions/r3x-docs/blob/master/README.md)

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
#### Supported Languages 
RubiX supports the following:

- JavaScript : `--type js`
- Golang : `--type go`
- Python : `--type py`

### Create a Function as a Container Image
The create function image locally:
```
$ cd <<function dir>>
$ r3x build -n <<Repo Name or Org>>
```
To create and push an image to Docker Hub:
```
$ r3x build -p -n <<Repo Name or Org>>

```
To create and push an image to Quay.io
```
$ r3x build -p -q -n <<Repo Name or Org>>
```

### Deploy a Function as a Container 
On success of the create command, a `service.yaml` file is generated and can be used to deploy the Function to Knative. For full instructions on this process see our [Documentation](https://github.com/rubixFunctions/r3x-docs/blob/master/install/README.md)

Alternatively, a Function as a Container can be run like any other container using a runtime like Docker, for example:
```
$ docker run -t -p 8080:8080 <<image tag>>
```


## License
This project is licensed under the Apache 2.0 License - see the [LICENSE](LICENSE) file for details

