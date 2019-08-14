// Tasker - A pluggable task server for keeping track of all those To-Do's
// Copyright (C) 2019 Ryan Murphy <ryan@oatmealrais.in>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"plugin"

	"github.com/oatmealraisin/tasker/pkg/plugins"
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

	loadedPlugins []plugins.TaskerPlugin
)

// TaskerCmd represents the base command when called without any subcommands
var TaskerCmd = &cobra.Command{
	Use:   "tasker",
	Short: "A pluggable task server for keeping track of all those To Do's",
	Long:  ``,
	Run:   statusCmd.Run,
	RunE:  status,

	//PreRun:  preRun,
	//PostRun: postRun,
	//TraverseChildren: false,
}

// Add persistent flags to the command, initialize the configuration.
func init() {
	cobra.OnInitialize(initConfig)

	TaskerCmd.PersistentFlags().StringVarP(&cfg, "config", "c", "", "Config directory to use.")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetDefault("StorageType", "sqlite")
	viper.SetDefault("WorkingDir", "$XDG_DATA_HOME/tasker")
	viper.SetDefault("Delimiter", "|")
	viper.SetDefault("EditCmd", "vi")
	viper.SetDefault("PluginDir", "$XDG_CONFIG_HOME/tasker/autoload")

	viper.SetEnvPrefix("tasker")
	// This means that any config variable can be set using the corresponding
	// TASKER_* environment variable
	viper.AutomaticEnv()

	if cfg != "" {
		f, err := os.Open(cfg)
		if err != nil {
			panic(err.Error())
		}
		defer f.Close()

		err = viper.ReadConfig(f)
		if err != nil {
			panic(err.Error())
		}
	} else {
		viper.AddConfigPath("$XDG_CONFIG_HOME/tasker/")
		viper.AddConfigPath("$HOME/.config/tasker/")
		viper.AddConfigPath("/etc/tasker/")
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")

		if err := viper.ReadInConfig(); err != nil {
			fmt.Printf("Could not read config file: %s\n", err.Error())
			os.Exit(1)
		}
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
		fmt.Fprintf(os.Stderr, "Unknown database type: %s", viper.GetString("StorageType"))
	}

	var err error
	termWidth, termHeight, err = terminal.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		termWidth = 200
	}

	pluginDir := viper.GetString("PluginDir")
	if pluginDir != "" {
		pluginDir = os.ExpandEnv(pluginDir)
		fns, err := ioutil.ReadDir(pluginDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: Could not open plugin directory %s: %s\n", pluginDir, err.Error())
		}

		for _, fn := range fns {
			//fmt.Printf("%s\n", filepath.Join(pluginDir, fn.Name()))
			p, err := plugin.Open(filepath.Join(pluginDir, fn.Name()))
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: Could not load plugin %s: %s\n", fn.Name(), err.Error())
			}

			sym, err := p.Lookup("Plugin")
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: Plugin '%s' does not have a Plugin object: %s\n", fn.Name(), err.Error())
			}

			newPlugin, ok := sym.(plugins.TaskerPlugin)
			if !ok {
				fmt.Fprintf(os.Stderr, "ERROR: Plugin '%s' does not implement the TaskerPlugin interface!\n", fn.Name())
			}

			loadedPlugins = append(loadedPlugins, newPlugin)

			go func() {
				newPlugin.Initialize()
			}()
		}
	}
}

// Execute adds all child commands to the root command and sets flags
// appropriately.  This is called by main.main(). It only needs to happen once
// to the TaskerCmd.
func Execute() error {
	return TaskerCmd.Execute()
}
