package usecase

import (
	"encoding/json"
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

func (r RouteBuilder) GetRouterGroup(item *openapi3.PathItem) string {
	group, ok := item.Extensions["x-go-gin-group"]

	if !ok {
		return "/"
	} else if v, ok := group.(json.RawMessage); !ok {
		return "/"
	} else {
		groupName := ""
		err := json.Unmarshal(v, &groupName)
		if err != nil {
			return "/"
		}
		if groupName == "" {
			return ""
		} else {
			return "/" + groupName
		}
	}
}

func (r RouteBuilder) BuildRouteSchema(pathPattern string, item *openapi3.PathItem) entity.Route {
	name := strcase.ToCamel(strings.Replace(pathPattern, "/", "_", -1))
	if name == "" {
		name = "Root"
	}
	route := entity.Route{
		Name:  name,
		Path:  pathPattern,
		Group: r.GetRouterGroup(item),
	}
	methods := r.BuildMethods(item)
	route.Methods = methods
	return route
}

func (r RouteBuilder) BuildRoutesSchema() []entity.Route {
	routes := []entity.Route{}
	groups := []string{}
	routeMap := map[string][]entity.Route{}
	for p, pi := range r.sp {
		routes = append(routes, r.BuildRouteSchema(p, pi))
	}

	for _, r := range routes {
		if routeMap[r.Group] == nil {
			groups = append(groups, r.Group)
			routeMap[r.Group] = []entity.Route{}
		}
		routeMap[r.Group] = append(routeMap[r.Group], r)
	}
	sort.Strings(groups)

	routes = []entity.Route{}
	for _, key := range groups {
		sort.SliceStable(routeMap[key], func(i int, j int) bool {
			if routeMap[key][i].Path < routeMap[key][j].Path {
				return true
			}
			return false
		})
		for _, r := range routeMap[key] {
			routes = append(routes, r)
		}
	}

	return routes
}

func NewRouteBuilder(sp openapi3.Paths, s SchemaBuilder) RouteBuilder {
	return RouteBuilder{
		sp: sp,
		sb: s,
	}
}
