package cmd

type Function struct {
	absPath string
	cmdPath string
	srcPath string
	//license License
	name string
}

// Returns Function with a specified function name
func NewFunction(functionName string) *Function {

}

// Returns Function with a specified absolute path to package
func NewFunctionFromPath(absPath string) *Function {

}
