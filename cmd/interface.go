package cmd

import (
	"github.com/spf13/cobra"
	"mk-oapigen-go/usecase"
)

var interfaceCmd = &cobra.Command{
	Use: "interface",
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
		gen := usecase.NewRouteGenerator(packageName, schemaPackage)
		gen.PrintRoute(routes)
		return nil
	},
}

func init() {
	interfaceCmd.Flags().StringP("openapi", "f", "./openapi.yml", "Input OpenAPI Specification file path [ex. docs/openapi.yml] (default openapi.yml")
	interfaceCmd.Flags().StringP("package", "p", "entity", "go mod module name [ex. sample api] (required)")
	interfaceCmd.Flags().StringP("schema_package", "s", "entity", "schema package")
	RootCmd.AddCommand(interfaceCmd)
}
