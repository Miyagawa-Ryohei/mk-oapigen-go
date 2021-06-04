package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"mk-oapigen-go/generator"
)

var serverCmd = &cobra.Command{
	Use: "server",
	RunE: func(cmd *cobra.Command, arg []string) error {
		fmt.Println("generate server code")

		specFile, err := cmd.Flags().GetString("openapi")
		if err != nil {
			return err
		}

		outputPackage, err := cmd.Flags().GetString("endpoint")
		if err != nil {
			return err
		}

		projectRoot, err := cmd.Flags().GetString("server")
		if err != nil {
			return err
		}

		routes, err := cmd.Flags().GetString("routes")
		if err != nil {
			return err
		}

		if err = generator.GenenateServerSideCode(specFile); err != nil {
			return err
		}
		fmt.Println(outputPackage)
		fmt.Println(projectRoot)
		fmt.Println(routes)
		return nil
	},
}

func init() {
	serverCmd.Flags().StringP("openapi", "f", "openapi.yml", "Input OpenAPI Specification file path [ex. docs/openapi.yml] (default openapi.yml")
	serverCmd.Flags().StringP("endpoint", "e", "src/gateway/endpoints", "Endpoint source files directory [ex. for/example/endpoints] (default src/gateway/endpoints)")
	serverCmd.Flags().StringP("server", "s", "src/infra/server.go", "Server Definition File [ ex. for/example/server.go ] (default src/infra/server.go)")
	serverCmd.Flags().StringP("routes", "r", "src/infra/routes", "Route Register files directory [ ex. for/example/register/route ] default (src/infra/routes)")
	RootCmd.AddCommand(serverCmd)
}
