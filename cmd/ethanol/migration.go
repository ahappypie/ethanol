package ethanol

import (
	"github.com/ahappypie/ethanol/internal/ethanol"
	"github.com/spf13/cobra"
)

var (
	Catalog string
	Schema  string
)

func init() {
	migrationCmd.AddCommand(migrationGenerateCmd, migrationRunCmd, migrationRevertCmd, migrationRedoCmd)

	rootCmd.AddCommand(migrationCmd)

	migrationCmd.PersistentFlags().StringVarP(&Catalog, "use-catalog", "c", "", "prepends USE CATALOG <input> to your queries. optional.")
	//viper.BindPFlag("environment.catalog", migrationCmd.PersistentFlags().Lookup("use-catalog"))
	//viper.BindEnv("use-catalog", "CATALOG")

	migrationCmd.PersistentFlags().StringVarP(&Schema, "use-schema", "s", "", "prepends USE SCHEMA <input> to your queries. optional.")
	//viper.BindPFlag("environment.schema", migrationCmd.PersistentFlags().Lookup("use-schema"))
	//viper.BindEnv("use-schema", "SCHEMA")
}

var migrationCmd = &cobra.Command{
	Use:   "migration",
	Short: "Execute and rollback migrations",
	Long:  `Execute pending migrations and rollback previous migrations on Unity Catalog.`,
}

var migrationGenerateCmd = &cobra.Command{
	Use:   "generate [migration_name]",
	Short: "Generate new migration in the specified directory",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ethanol.GenerateMigration(Directory, args[0])
	},
}
var migrationRunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run all pending migrations",
	Run: func(cmd *cobra.Command, args []string) {
		ethanol.RunMigration(Directory)
	},
}

var migrationRevertCmd = &cobra.Command{
	Use:   "revert",
	Short: "Revert last run migration",
	Long:  "Revert last run migration as determined by the internal tracking table",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var migrationRedoCmd = &cobra.Command{
	Use:   "redo",
	Short: "Redo last migration",
	Long:  "Revert and re-run last migration as determined by the internal tracking table",
	Run: func(cmd *cobra.Command, args []string) {

	},
}
