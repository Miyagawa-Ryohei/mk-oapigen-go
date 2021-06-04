package types

import "fmt"

type RootDefinitions struct {
	Routes *[]RouteDefinition
	Server *ServerDefinition
	Types  *[]TypeDefinition
}

type RouteDefinition struct {
	FileName  string
	BaseName  string
	Path      string
	Endpoints []EndpointDefinition
}

type ServerDefinition struct {
	APIBasePath string
	Title       string
}

type EndpointDefinition struct {
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
