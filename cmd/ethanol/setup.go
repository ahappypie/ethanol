package ethanol

import (
	"fmt"
	"github.com/ahappypie/ethanol/internal/ethanol"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	rootCmd.AddCommand(setupCmd)
}

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Generate bootstrap SQL",
	Long:  "Generate bootstrap SQL including Catalog, Schema and Table setup, along with permissions and initial data.",
	Run: func(cmd *cobra.Command, args []string) {
		if InternalCatalog == "" {
			fmt.Println("internal catalog must be specified for this operation")
			os.Exit(1)
		}
		ethanol.GenerateSetup(Directory, InternalCatalog, InternalSchema, InternalTable)
	},
}
