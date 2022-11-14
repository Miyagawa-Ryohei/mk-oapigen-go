package cmd

import (
	"github.com/spf13/cobra"
	"mk-oapigen-go/usecase"
)

var serverCmd = &cobra.Command{
	Use: "server",
	RunE: func(cmd *cobra.Command, arg []string) error {

		specFile, err := cmd.Flags().GetString("openapi")
		if err != nil {
			return err
		}

		packageName, err := cmd.Flags().GetString("package")
		if err != nil {
			return err
		}

		routerPackage, err := cmd.Flags().GetString("router_package")
		if err != nil {
			return err
		}

		spec, err := usecase.ReadSpec(specFile)
		if err != nil {
			return err
		}
		serverBuilder := usecase.ServerBuilder{}
		schemaBuilder := usecase.SchemaBuilder{}
		routeBuilder := usecase.NewRouteBuilder(spec.Paths, schemaBuilder)
		server := serverBuilder.BuildServerSchema(spec.Servers[0])
		routes := routeBuilder.BuildRoutesSchema()

		if routerPackage == packageName {
			routerPackage = ""
		}
		gen := usecase.NewServerGenerator(packageName, routerPackage)
		gen.PrintServer(server, routes)

		return nil
	},
}

func init() {
	serverCmd.Flags().StringP("openapi", "f", "openapi.yml", "Input OpenAPI Specification file path [ex. docs/openapi.yml] (default openapi.yml")
	serverCmd.Flags().StringP("package", "p", "infra", "go mod module name [ex. sample api] (required)")
	serverCmd.Flags().StringP("router_package", "s", "entity", "schema package")
	RootCmd.AddCommand(serverCmd)
}
