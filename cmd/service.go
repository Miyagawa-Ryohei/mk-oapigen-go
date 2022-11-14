package cmd

import (
	"github.com/spf13/cobra"
	"mk-oapigen-go/usecase"
)

var serviceCmd = &cobra.Command{
	Use: "service",
	RunE: func(cmd *cobra.Command, arg []string) error {

		specFile, err := cmd.Flags().GetString("openapi")
		if err != nil {
			return err
		}

		packageName, err := cmd.Flags().GetString("package")
		if err != nil {
			return err
		}

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
		gen := usecase.NewServiceGenerator(packageName, schemaPackage)
		gen.PrintService(routes)

		return nil
	},
}

func init() {
	serviceCmd.Flags().StringP("openapi", "f", "openapi.yml", "Input OpenAPI Specification file path [ex. docs/openapi.yml] (default openapi.yml")
	serviceCmd.Flags().StringP("package", "p", "infra", "go mod module name [ex. sample api] (required)")
	serviceCmd.Flags().StringP("schema_package", "s", "entity", "schema package")
	RootCmd.AddCommand(serviceCmd)
}
