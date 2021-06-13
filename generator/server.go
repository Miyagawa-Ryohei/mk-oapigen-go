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
	"strconv"
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
		Types:  createTypeStruct(sp),
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

	root.Server = &server
	root.Routes = &routes

	return root, nil
}

func GenFromTypes(tmpl []byte, definitions map[string][]types.TypeDefinition, filename string) {
	writeBufToGoFile(tmpl, filename, definitions)
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

func addType(result map[string][]types.TypeDefinition, typeName string, ref *openapi3.Schema) {
	for k,content := range(ref.Properties) {
		propertyName := strcase.ToCamel(k)
		if content.Ref != "" {
			schema := strings.Split(content.Ref, "/")
			name := schema[len(schema)-1]
			result[name] = []types.TypeDefinition{}
			addType(result, name, content.Value)
			result[typeName] = append(result[typeName], types.TypeDefinition{
				Name:     propertyName,
				TypeName: name,
			})
			continue
		}
		if content.Value.Type == "object" {
			name := strcase.ToCamel(typeName + propertyName)
			addType(result, name, content.Value)
			result[typeName] = append(result[typeName], types.TypeDefinition{
				Name:     propertyName,
				TypeName: name,
			})
			continue
		}
		if content.Value.Type == "number" {
			result[typeName] = append(result[typeName], types.TypeDefinition{
				Name:     propertyName,
				TypeName: "float64",
			})
			continue
		}
		if content.Value.Type == "integer" {
			result[typeName] = append(result[typeName], types.TypeDefinition{
				Name:     propertyName,
				TypeName: "int64",
			})
			continue
		}
		if content.Value.Type == "boolean" {
			result[typeName] = append(result[typeName], types.TypeDefinition{
				Name:     propertyName,
				TypeName: "bool",
			})
			continue
		}
		if content.Value.Type == "string" {
			switch content.Value.Format {

			case "date":
			case "date-time":
				result[typeName] = append(result[typeName], types.TypeDefinition{
					Name:     propertyName,
					TypeName: "time.Time",
				})
				break;

			case "binary":
				result[typeName] = append(result[typeName], types.TypeDefinition{
					Name:     propertyName,
					TypeName: "[]byte",
				})
				break;

			default:
				result[typeName] = append(result[typeName], types.TypeDefinition{
					Name:     propertyName,
					TypeName: "string",
				})
				break;
			}
		}
	}

}

func createTypeStruct(oapiDef *openapi3.T) map[string][]types.TypeDefinition {
	result := make(map[string][]types.TypeDefinition)

	for _, v := range(oapiDef.Paths) {
		for _, methodContents := range v.Operations() {
			bass := strcase.ToCamel(methodContents.OperationID)
			fmt.Println(bass)
			for code, res  := range methodContents.Responses {
				responseTypeName := ""
				statusCode, err := strconv.Atoi(code)

				if err != nil {
					panic(err)
				}

				if statusCode - 200 < 100 {
					responseTypeName = bass+"Response"
				} else {
					responseTypeName = bass + code + "Error"
				}

				result[responseTypeName] = []types.TypeDefinition{}

				for mime, content := range res.Value.Content {
					if mime == "application/json" {
						if content.Schema.Ref != "" {
							schema := strings.Split(content.Schema.Ref,"/")
							name := schema[len(schema)-1]
							addType(result, name, content.Schema.Value)
						}
					} else {
						result[responseTypeName] = append(result[responseTypeName],types.TypeDefinition{
							Name: "",
							TypeName: "string",
						})
					}
				}
			}

			requestTypeName :=  bass+"Request"
			result[requestTypeName] = []types.TypeDefinition{}
			if methodContents.RequestBody == nil {
				continue
			}
			for mime, content := range methodContents.RequestBody.Value.Content{
				if mime == "application/json" {
					if content.Schema.Ref != "" {
						schema := strings.Split(content.Schema.Ref,"/")
						name := schema[len(schema)-1]
						addType(result, name, content.Schema.Value)
					}
				} else {
					result[requestTypeName] = append(result[requestTypeName],types.TypeDefinition{
						Name: "",
						TypeName: "string",
					})
				}
			}
		}

	}

	return result
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
	GenFromTypes(tempBuf["template/server/types.gotmpl"], root.Types, path.Join(src,"entity","types"))
	GenFromRouteDefinition(tempBuf["template/server/route.gotmpl"], root.Routes, path.Join(src,"infra","registry","route"))
	GenFromRouteDefinition(tempBuf["template/server/endpoint.gotmpl"], root.Routes, path.Join(src,p))
	GenFromStatic(tempBuf["template/server/application.gotmpl"], path.Join(src, "infra", "registry","initialize"))
	GenFromServer(tempBuf["template/server/server.gotmpl"], root.Server, path.Join(src, "infra", "registry", "server"))
	GenFromServer(tempBuf["template/server/root_cmd.gotmpl"], root.Server, path.Join(src, "adapter", "cmd", "root_cmd"))
	GenFromRoot(tempBuf["template/server/server_cmd.gotmpl"], root, path.Join(src, "adapter", "cmd","server_cmd"))

	return nil
}
