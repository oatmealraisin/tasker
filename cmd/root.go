package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfg string

// TaskerCmd represents the base command when called without any subcommands
var TaskerCmd = &cobra.Command{
	Use:   "tasker",
	Short: "A pluggable task server for keeping track of all those To Do's",
	Long:  ``,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the TaskerCmd.
func Execute() error {
	return TaskerCmd.Execute()
}

// Add persistent flags to the command, initialize the configuration.
func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	TaskerCmd.PersistentFlags().StringVarP(&cfg, "config", "c", "", "Config directory to use.")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// The config name is always "config". May make this changeable
	viper.SetConfigName("config")
	if cfg != "" {
		viper.AddConfigPath(cfg)
	} else {
		viper.AddConfigPath("$XDG_CONFIG_HOME/tasker")
		viper.AddConfigPath("$HOME/.tasker")
		viper.AddConfigPath("/etc/tasker")
	}
}
