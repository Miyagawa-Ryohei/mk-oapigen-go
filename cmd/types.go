package cmd

import (
	"github.com/spf13/cobra"
	"mk-oapigen-go/usecase"
)

var typeCmd = &cobra.Command{
	Use: "type",
	RunE: func(cmd *cobra.Command, arg []string) error {

		specFile, err := cmd.Flags().GetString("openapi")
		if err != nil {
			return err
		}

		packageName, err := cmd.Flags().GetString("package")
		if err != nil {
			return err
		}
		spec, err := usecase.ReadSpec(specFile)
		if err != nil {
			return err
		}
		schemaBuilder := usecase.SchemaBuilder{}
		routeBuilder := usecase.NewRouteBuilder(spec.Paths, schemaBuilder)
		routes := routeBuilder.BuildRoutesSchema()
		gen := usecase.NewTypeGenerator(packageName)
		schemas := gen.ExtractTypes(routes)
		gen.PrintType(schemas)
		return nil
	},
}

func init() {
	typeCmd.Flags().StringP("openapi", "f", "openapi.yml", "Input OpenAPI Specification file path [ex. docs/openapi.yml] (default openapi.yml")
	typeCmd.Flags().StringP("package", "p", "", "go mod module name [ex. sample api] (required)")
	RootCmd.AddCommand(typeCmd)
}
