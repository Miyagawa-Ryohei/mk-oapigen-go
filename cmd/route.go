package cmd

import (
	"github.com/spf13/cobra"
	"mk-oapigen-go/usecase"
)

var routerCmd = &cobra.Command{
	Use: "router",
	RunE: func(cmd *cobra.Command, arg []string) error {

		specFile, err := cmd.Flags().GetString("openapi")
		if err != nil {
			return err
		}

		packageName, err := cmd.Flags().GetString("package")

		schemaPackage, err := cmd.Flags().GetString("schema_package")
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

		if schemaPackage == packageName {
			schemaPackage = ""
		}
		gen := usecase.NewRouterGenerator(packageName, schemaPackage)
		gen.PrintRouter(routes)
		return nil
	},
}

func init() {
	routerCmd.Flags().StringP("openapi", "f", "openapi.yml", "Input OpenAPI Specification file path [ex. docs/openapi.yml] (default openapi.yml")
	routerCmd.Flags().StringP("package", "p", "", "go mod module name [ex. sample api] (required)")
	routerCmd.Flags().StringP("schema_package", "s", "entity", "schema package")
	RootCmd.AddCommand(routerCmd)
}
