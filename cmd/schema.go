package cmd

import "fmt"

type Schema struct {
	name     string
	funcType string
	response string
}

// NewSchema Returns Schema with a specified function name
func NewSchema(functionName string, funcType string, response string) *Schema {
	if functionName == "" {
		fmt.Println("can't create function with no name")
	}
	if funcType == "" {
		fmt.Println("can't create function with no type")
	}
	s := new(Schema)

	if response == "" {
		s.response = "JSON"
	}

	s.response = response
	s.name = functionName
	s.funcType = funcType

	return s
}
