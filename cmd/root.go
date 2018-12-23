package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/oatmealraisin/tasker/pkg/storage"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"golang.org/x/crypto/ssh/terminal"
)

var (
	cfg        string
	db         storage.Storage
	termWidth  int
	termHeight int
)

// TaskerCmd represents the base command when called without any subcommands
var TaskerCmd = &cobra.Command{
	Use:   "tasker",
	Short: "A pluggable task server for keeping track of all those To Do's",
	Long:  ``,
}

// Execute adds all child commands to the root command and sets flags
// appropriately.  This is called by main.main(). It only needs to happen once
// to the TaskerCmd.
func Execute() error {
	return TaskerCmd.Execute()
}

// Add persistent flags to the command, initialize the configuration.
func init() {
	cobra.OnInitialize(initConfig)

	TaskerCmd.PersistentFlags().StringVarP(&cfg, "config", "c", "", "Config directory to use.")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// The config name is always "config". May make this changeable
	if cfg != "" {
		viper.SetConfigFile(cfg)
	} else {
		viper.AddConfigPath("$XDG_CONFIG_HOME/tasker/")
		viper.AddConfigPath("$HOME/.config/tasker/")
		viper.AddConfigPath("/etc/tasker/")
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.SetDefault("StorageType", "sqlite")
	viper.SetDefault("WorkingDir", "$XDG_DATA_HOME/tasker")
	viper.SetDefault("Delimiter", "|")
	viper.SetDefault("EditCmd", "vi")

	viper.SetEnvPrefix("tasker")
	viper.AutomaticEnv()
	// This means that any config variable can be set using the corresponding
	// TASKER_* environment variable

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Could not read config file: %s\n", err.Error())
		os.Exit(1)
	}

	switch viper.GetString("StorageType") {
	case "csv":
		db = storage.NewCsvStorage(filepath.Join(viper.GetString("WorkingDir"), "tasklist.csv"))
		break
	case "sqlite":
		panic("sqlite not implemented")
	case "postgres":
		panic("postgres not implemented")
	default:
		fmt.Printf("Unknown database type: %s", viper.GetString("StorageType"))
	}

	var err error
	termWidth, termHeight, err = terminal.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		panic(err.Error())
	}
}
