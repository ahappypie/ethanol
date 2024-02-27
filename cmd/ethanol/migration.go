package ethanol

import (
	"github.com/ahappypie/ethanol/internal/ethanol"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	Catalog      string
	Schema       string
	ClientId     string
	ClientSecret string
	Host         string
	HttpPath     string
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

	viper.BindEnv("cluster.clientId", "DATABRICKS_CLIENT_ID")
	viper.BindEnv("cluster.clientSecret", "DATABRICKS_CLIENT_SECRET")
	viper.BindEnv("cluster.host", "DATABRICKS_HOST")
	viper.BindEnv("cluster.httpPath", "DATABRICKS_HTTP_PATH")
}

func initMigrationConfig() {
	//unmarshal
	viper.UnmarshalKey("cluster.clientId", &ClientId)
	viper.UnmarshalKey("cluster.clientSecret", &ClientSecret)
	//TODO mark these as required
	viper.UnmarshalKey("cluster.host", &Host)
	viper.UnmarshalKey("cluster.httpPath", &HttpPath)
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
	PreRun: func(cmd *cobra.Command, args []string) {
		initMigrationConfig()
	},
	Run: func(cmd *cobra.Command, args []string) {
		ethanol.RunMigration(Directory, InternalCatalog, InternalSchema, InternalTable, ClientId, ClientSecret, Host, HttpPath)
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
