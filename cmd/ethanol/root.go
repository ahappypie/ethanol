package ethanol

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var (
	UseConfig       bool
	ConfigFile      string
	Directory       string
	InternalCatalog string
	InternalSchema  string
	InternalTable   string
)

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().BoolVar(&UseConfig, "use-config", false, "should use the yaml config. optional.")
	rootCmd.PersistentFlags().StringVar(&ConfigFile, "config", "ethanol.yaml", "config file. unused unless --use-config is present. optional.")

	rootCmd.PersistentFlags().StringP("directory", "d", "migrations", "relative path to migrations, i.e. ../repository/migrations. optional.")
	viper.BindPFlag("directory", rootCmd.PersistentFlags().Lookup("directory"))
	viper.BindEnv("directory", "DIRECTORY")

	rootCmd.PersistentFlags().StringP("internal-catalog", "i", "", "catalog where the internal tracking table is stored, required.")
	viper.BindPFlag("internal.catalog", rootCmd.PersistentFlags().Lookup("internal-catalog"))
	viper.BindEnv("internal.catalog, INTERNAL_CATALOG")

	rootCmd.PersistentFlags().String("internal-schema", "default", "schema where the internal tracking table is stored, optional.")
	viper.BindPFlag("internal.schema", rootCmd.PersistentFlags().Lookup("internal-schema"))
	viper.BindEnv("internal.schema", "INTERNAL_SCHEMA")

	rootCmd.PersistentFlags().String("internal-table", "ethanol_migrations", "table where the internal tracking data is stored, optional.")
	viper.BindPFlag("internal.table", rootCmd.PersistentFlags().Lookup("internal-table"))
	viper.BindEnv("internal.table", "INTERNAL_TABLE")
}

func initConfig() {
	if !UseConfig {
		viper.AutomaticEnv()
	} else {
		viper.SetConfigFile(ConfigFile)
		if err := viper.ReadInConfig(); err != nil {
			fmt.Println("Can't read config:", err)
		}
	}

	//unmarshal
	viper.UnmarshalKey("directory", &Directory)
	if err := viper.UnmarshalKey("internal.catalog", &InternalCatalog); err != nil {
		fmt.Println("error reading internal catalog")
		os.Exit(1)
	}
	viper.UnmarshalKey("internal.schema", &InternalSchema)
	viper.UnmarshalKey("internal.table", &InternalTable)
}

var rootCmd = &cobra.Command{
	Use:   "ethanol",
	Short: "Runs SQL migrations on Unity Catalog",
	Long: `Ethanol runs SQL migrations on Unity Catalog,
			enabling CI workflows on your data warehouse.`,
}
