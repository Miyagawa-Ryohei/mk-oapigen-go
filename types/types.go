package types

type RootDefinitions struct {
	Module         string
	EndpointPackage string
	Routes *[]RouteDefinition
	Server *ServerDefinition
	Types   map[string][]TypeDefinition
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
	TypeName string
}
