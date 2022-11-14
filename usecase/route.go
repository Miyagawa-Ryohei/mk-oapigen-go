package usecase

import (
	"bytes"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"go/format"
	"mk-oapigen-go/entity"
)

type RouteGenerator struct {
	Package       string
	SchemaPackage string
}

func (r RouteGenerator) GenerateRoute(spec openapi3.T) error {
	return nil
}
func (r RouteGenerator) CreateHeader() string {
	return fmt.Sprintf("// Package %s \n"+
		"//this file is generated by swag2go. DO NOT EDIT\n"+
		"//    authorized : miyagawa.ryohei\n"+
		"package %s\n"+
		"\n"+
		"import ( \n"+
		"\"github.com/gin-gonic/gin\"\n"+
		")\n", r.Package, r.Package)
}
func (r RouteGenerator) CreateMethodList(methods []entity.Method) (str string) {
	for _, m := range methods {
		typePackage := fmt.Sprintf("%s.", r.SchemaPackage)
		if r.SchemaPackage == "" {
			typePackage = ""
		}
		arg := fmt.Sprintf("req %s%s", typePackage, m.Request.Name)
		if m.Request.Name == "" {
			arg = ""
		}
		str = str + fmt.Sprintf(
			"//%s HTTP %s Handler\n"+
				"%s(%s) (*%s%s, error)\n\n",
			m.Name,
			m.Type,
			m.Name,
			arg,
			typePackage,
			m.Response.Name,
		)
	}
	return
}

func (r RouteGenerator) CreateRouteInterface(route entity.Route) (str string) {
	str = fmt.Sprintf(
		"//%sService \n"+
			"// path : %s \n"+
			"type %sService interface { \n\n %s }\n\n",
		route.Name,
		route.Path,
		route.Name,
		r.CreateMethodList(route.Methods),
	)
	return
}

func (r RouteGenerator) CreateRouteProvider(route entity.Route) (str string) {
	str = fmt.Sprintf(
		"type %sServiceProvider func() (%sService, error)\n\n",
		route.Name,
		route.Name,
	)
	return
}

func (r RouteGenerator) CreateHandlerFuncProvider() (str string) {
	str = fmt.Sprintf(
		"type HandlerFuncProvider func() gin.HandlerFunc\n\n",
	)
	return
}

func (r RouteGenerator) CreateRoute(routes []entity.Route) (str string) {
	for _, route := range routes {
		str = str + r.CreateRouteInterface(route)
		str = str + r.CreateRouteProvider(route)
	}
	str = str + r.CreateHandlerFuncProvider()
	return
}

func (r RouteGenerator) PrintRoute(routes []entity.Route) (str string) {
	str = fmt.Sprintf("%s %s", r.CreateHeader(), r.CreateRoute(routes))

	written, err := format.Source(bytes.NewBufferString(str).Bytes())
	if err != nil {
		panic(err)
	}
	fmt.Println(string(written))
	return
}

func NewRouteGenerator(p string, schemaPackage string) RouteGenerator {
	return RouteGenerator{
		Package:       p,
		SchemaPackage: schemaPackage,
	}
}
