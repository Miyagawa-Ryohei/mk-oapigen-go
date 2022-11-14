package usecase

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/iancoleman/strcase"
	"mk-oapigen-go/entity"
	"sort"
	"strings"
)

type RouteBuilder struct {
	sp openapi3.Paths
	sb SchemaBuilder
}

func (r RouteBuilder) BuildMethod(method string, operation *openapi3.Operation) entity.Method {
	m := entity.Method{
		Type:        strings.ToUpper(method),
		Name:        operation.OperationID,
		Description: operation.Description,
		Summary:     operation.Summary,
	}
	requestParam := entity.Schema{}
	if operation.RequestBody != nil {
		requestParam = *(r.sb.BuildSchemaFromRequestBody(operation.RequestBody))
	} else if len(operation.Parameters) > 0 {
		requestParam = r.sb.BuildSchemaFromParameters(&operation.Parameters, operation.OperationID)
	} else {
		requestParam = entity.Schema{}
	}

	response := entity.Schema{}
	if respRef := operation.Responses.Get(200); respRef != nil {
		response = *(r.sb.BuildSchemaFromRespRef(respRef))
	} else {
		response = entity.Schema{}
	}
	m.Request = requestParam
	m.Response = response
	return m
}

func (r RouteBuilder) BuildMethods(item *openapi3.PathItem) []entity.Method {
	ret := []entity.Method{}
	for method, op := range item.Operations() {
		ret = append(ret, r.BuildMethod(method, op))
	}
	return ret
}

func (r RouteBuilder) BuildRouteSchema(pathPattern string, item *openapi3.PathItem) entity.Route {
	name := strcase.ToCamel(strings.Replace(pathPattern, "/", "_", -1))
	if name == "" {
		name = "Root"
	}
	route := entity.Route{
		Name: name,
		Path: pathPattern,
	}
	methods := r.BuildMethods(item)
	route.Methods = methods
	return route
}

func (r RouteBuilder) BuildRoutesSchema() []entity.Route {
	routes := []entity.Route{}
	for p, pi := range r.sp {
		routes = append(routes, r.BuildRouteSchema(p, pi))
	}
	sort.SliceStable(routes, func(i int, j int) bool {
		if routes[i].Path < routes[j].Path {
			return true
		}
		return false
	})
	return routes
}

func NewRouteBuilder(sp openapi3.Paths, s SchemaBuilder) RouteBuilder {
	return RouteBuilder{
		sp: sp,
		sb: s,
	}
}
