package generator

import (
	"bytes"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/iancoleman/strcase"
	"go/format"
	"mk-oapigen-go/types"
	"os"
	"path"
	"strings"
	"text/template"
)

func getServerTempNames() ([]string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	dirs, err := os.ReadDir(path.Join(cwd, "template", "server"))
	if err != nil {
		return nil, err
	}

	var templates []string
	for _, fi := range dirs {
		if fi.IsDir() {
			continue
		}
		templates = append(templates, path.Join("template", "server", fi.Name()))
	}

	return templates, nil
}

func loadFromSpec(openapi string, module string) (*types.RootDefinitions, error) {
	sp, err := openapi3.NewLoader().LoadFromFile(openapi)
	if err != nil {
		return nil, err
	}
	root := &types.RootDefinitions{
		Module: module,
		EndpointPackage: path.Join(module,"adapter","gateway","api",sp.Info.Version),
		Routes: nil,
		Server: nil,
		Types:  nil,
	}

	server := types.ServerDefinition{
		APIBasePath: sp.Servers[0].URL,
		Title:       sp.Info.Title,
		Module: module,
	}

	var routes []types.RouteDefinition

	for pathStr, pathContents := range sp.Paths {
		FileName := pathContents.Summary
		BaseName := strcase.ToCamel(pathContents.Summary)
		var Endpoints []types.EndpointDefinition
		for method, methodContents := range pathContents.Operations() {
			Endpoints = append(Endpoints, types.EndpointDefinition{
				Module: module,
				Path:         pathStr,
				Method:       method,
				FunctionName: strcase.ToCamel(methodContents.OperationID),
			})
		}
		routes = append(routes, types.RouteDefinition{
			Module: module,
			FileName:  FileName,
			BaseName:  BaseName,
			Path:      pathStr,
			Endpoints: Endpoints,
		})
	}

	var Types []types.TypeDefinition
	root.Server = &server
	root.Routes = &routes
	root.Types = &Types

	return root, nil
}

func GenFromServer(tmpl []byte, definitions *types.ServerDefinition, filename string) {
	writeBufToGoFile(tmpl, filename, definitions)
}

func GenFromRoot(tmpl []byte, definitions *types.RootDefinitions, filename string) {
	writeBufToGoFile(tmpl, filename, definitions)
}

func GenFromStatic(tmpl []byte, filename string) {
	writeBufToGoFile(tmpl, filename, nil)
}

func GenFromRouteDefinition(tmpl []byte, definitions *[]types.RouteDefinition, src string) {
	for _, definition := range *definitions {
		writeBufToGoFile(tmpl, path.Join(src, definition.FileName), definition)
	}
}

func writeBufToGoFile(tmpl []byte, filebase string, param interface{}) {
	b := &bytes.Buffer{}

	if err := template.Must(template.New("call").Parse(string(tmpl))).Execute(b, param); err != nil {
		panic(err)
	}
	fmt.Println(string(b.Bytes()))

	written, err := format.Source(b.Bytes())
	if err != nil {
		panic(err)
	}

	f, err := os.Create((filebase) + ".go")

	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	if err != nil {
		panic(err)
	}
	if _, err := f.Write(written); err != nil {
		panic(err)
	}
}

func GenerateServerSideCode(specFile string, module string, src string) error {

	templates, err := getServerTempNames()
	if err != nil {
		return err
	}

	tempBuf := map[string][]byte{}
	for _, temp := range templates {
		f, err := os.ReadFile(temp)
		if err != nil {
			return err
		}
		tempBuf[temp] = f
	}

	root, err := loadFromSpec(specFile, module)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(path.Join(src, "infra", "registry", "route"),509) ; err != nil {
		return err
	}
	if err := os.MkdirAll(path.Join(src, "adapter", "cmd"),509) ; err != nil {
		return err
	}
	if err := os.MkdirAll(path.Join(src, "adapter", "cmd"),509) ; err != nil {
		return err
	}
	ps := strings.Split(root.EndpointPackage,"/")
	p := path.Join(ps[1:]...)
	if err := os.MkdirAll(path.Join(src, p),509) ; err != nil {
		return err
	}
	fmt.Println(p)
	fmt.Println(src)
	GenFromRouteDefinition(tempBuf["template/server/route.gotmpl"], root.Routes, path.Join(src,"infra","registry","route"))
	GenFromRouteDefinition(tempBuf["template/server/endpoint.gotmpl"], root.Routes, path.Join(src,p))
	GenFromStatic(tempBuf["template/server/application.gotmpl"], path.Join(src, "infra", "registry","initialize"))
	GenFromServer(tempBuf["template/server/server.gotmpl"], root.Server, path.Join(src, "infra", "registry", "server"))
	GenFromServer(tempBuf["template/server/root_cmd.gotmpl"], root.Server, path.Join(src, "adapter", "cmd", "root_cmd"))
	GenFromRoot(tempBuf["template/server/server_cmd.gotmpl"], root, path.Join(src, "adapter", "cmd","server_cmd"))

	return nil
}
