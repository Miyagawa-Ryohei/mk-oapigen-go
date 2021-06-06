package types

import "fmt"

type RootDefinitions struct {
	Module         string
	EndpointPackage string
	Routes *[]RouteDefinition
	Server *ServerDefinition
	Types  *[]TypeDefinition
}

type RouteDefinition struct {
	Module         string
	FileName  string
	BaseName  string
	Path      string
	Endpoints []EndpointDefinition
}

type ServerDefinition struct {
	Module         string
	APIBasePath string
	Title       string
}

type EndpointDefinition struct {
	Module         string
	Path         string
	Method       string
	FunctionName string
}

type TypeDefinition struct {
	Name     string
	Structed string
}

func (t TypeDefinition) toString() string {
	return fmt.Sprintf("%s    %s", t.Name, t.Structed)
}
