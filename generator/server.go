package generator

import (
	"bytes"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/iancoleman/strcase"
	"go/format"
	"os"
	"path"
	"swag2go/types"
	"text/template"
)

func getServerTempNames () ([]string, error){
	cwd, err := os.Getwd()
	if err!=nil {
		return nil, err
	}

	dirs, err := os.ReadDir(path.Join(cwd,"template", "server"))
	if err!=nil {
		return nil, err
	}

	var templates = []string{}
	for _, fi := range dirs {
		if fi.IsDir() {
			continue
		}
		templates = append(templates, path.Join("template", "server", fi.Name()))
	}

	return templates, nil
}

func loadFromSpec(openapi string) (*types.RootDefinitions, error){
	root := &types.RootDefinitions{
		Routes: nil,
		Server: nil,
		Types:  nil,
	}
	sp, err := openapi3.NewLoader().LoadFromFile(openapi)
	if err != nil {
		return nil, err
	}
	server := types.ServerDefinition{
		APIBasePath: sp.Servers[0].URL,
		Title: sp.Info.Title,
	}
	routes := []types.RouteDefinition{}
	for path, pathContents := range sp.Paths {
		FileName := pathContents.Summary
		BaseName := strcase.ToCamel(pathContents.Summary)
		Endpoints := []types.EndpointDefinition{}
		for method, method_contents := range pathContents.Operations(){
			Endpoints = append(Endpoints, types.EndpointDefinition{
				Path:         path,
				Method:       method,
				FunctionName: strcase.ToCamel(method_contents.OperationID),
			})
		}
		routes = append(routes, types.RouteDefinition{
			FileName: FileName,
			BaseName:  BaseName,
			Path:      path,
			Endpoints: Endpoints,
		})
	}

	Types := []types.TypeDefinition{}
	root.Server = &server
	root.Routes = &routes
	root.Types = &Types

	return root, nil
}

func GenFromServer(tmpl []byte, definitions *types.ServerDefinition,filename string) {
	writeBufToGoFile(tmpl, filename, definitions)
}

func GenFromRoot(tmpl []byte, definitions *types.RootDefinitions, filename string) {
	writeBufToGoFile(tmpl, filename, definitions)
}

func GenFromStatic(tmpl []byte, filename string) {
	writeBufToGoFile(tmpl, filename , nil)
}

func GenFromRouteDefinition(tmpl []byte, definitions *[]types.RouteDefinition) {
	for _, definition := range *definitions {
		writeBufToGoFile(tmpl, definition.FileName, definition)
	}
}

func writeBufToGoFile(tmpl []byte, filebase string, param interface{}) {
	b := &bytes.Buffer{}

	if err := template.Must(template.New("call").Parse(string(tmpl))).Execute(b, param); err != nil {
		panic(err)
	}
	fmt.Println("%s", string(b.Bytes()))
	written, err := format.Source(b.Bytes())
	if err != nil {
		panic(err)
	}
	f,err := os.Create((filebase)+".go")
	defer f.Close()
	if err != nil {
		panic(err)
	}
	if _, err := f.Write(written); err != nil {
		panic(err)
	}
}

func GenenateServerSideCode(openapi string) error{

	templates,err := getServerTempNames()
	if err != nil {
		return err
	}

	tempBuf := map[string][]byte{}
	for _, template := range templates{
		f, err := os.ReadFile(template)
		if err != nil {
			return err
		}
		tempBuf[template] = f
	}

	root, err := loadFromSpec(openapi)
	if err != nil {
		return err
	}

	GenFromRouteDefinition(tempBuf["template/server/route.gotmpl"], root.Routes)
	GenFromRouteDefinition(tempBuf["template/server/endpoint.gotmpl"], root.Routes)
	GenFromStatic(tempBuf["template/server/application.gotmpl"],"initialize")
	GenFromServer(tempBuf["template/server/server.gotmpl"], root.Server,"server")
	GenFromServer(tempBuf["template/server/root_cmd.gotmpl"], root.Server, "root_cmd")
	GenFromRoot(tempBuf["template/server/server_cmd.gotmpl"], root, "server_cmd")
	if err != nil {
		return err
	}

	return nil
}
