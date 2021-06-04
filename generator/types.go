package generator

type RootDefinitions struct {
	Routes *RouteDefinitions
	Server *ServerDefinition
	Types  *TypeDefinitions
}

type RouteDefinitions = []RouteDefinition

type RouteDefinition struct {
	BaseName string
	Path  string
	EndpointDefinitions
}

type ServerDefinition struct {
	APIBasePath string
}

type EndpointDefinitions = []EndpointDefinition

type EndpointDefinition struct {
	Path  string
	Method string
	FunctionName string
}


type TypeDefinitions = []TypeDefinition

type TypeDefinition struct {
	structed string
}