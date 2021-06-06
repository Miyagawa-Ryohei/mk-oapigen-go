package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"mk-oapigen-go/generator"
	"os"
	"path"
)

var serverCmd = &cobra.Command{
	Use: "server",
	RunE: func(cmd *cobra.Command, arg []string) error {
		fmt.Println("generate server code")

		specFile, err := cmd.Flags().GetString("openapi")
		if err != nil {
			return err
		}

		module, err := cmd.Flags().GetString("module")
		if err != nil {
			return err
		}

		src, err := cmd.Flags().GetString("src")
		if err != nil {
			return err
		}

		root, err := cmd.Flags().GetString("root")
		if err != nil {
			return err
		}

		srcDir := path.Join(root,src)

		if err := os.MkdirAll(path.Join(srcDir),509) ; err != nil {
			return err
		}
		fmt.Println(srcDir)
		if err = generator.GenerateServerSideCode(specFile,module,srcDir); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	serverCmd.Flags().StringP("openapi", "f", "openapi.yml", "Input OpenAPI Specification file path [ex. docs/openapi.yml] (default openapi.yml")
	serverCmd.Flags().StringP("module", "m", "", "go mod module name [ex. sample api] (required)")
	serverCmd.Flags().StringP("src", "s", "./", "source directory [ex. src/] (default src ./)")
	serverCmd.Flags().StringP("root", "r", "./", "project root directory [ex. ./project] (default src ./)")
	RootCmd.AddCommand(serverCmd)
}
