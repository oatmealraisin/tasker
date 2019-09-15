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
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/oatmealraisin/tasker/cmd"
	"github.com/oatmealraisin/tasker/pkg/plugins"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	pflag.CommandLine.ParseErrorsWhitelist.UnknownFlags = true
	var config *string = pflag.StringP("config", "c", "", "")
	var noPlugins *bool = pflag.Bool("no-plugins", false, "")
	// Make sure pflag doesn't grab the program when we ask for help..
	var _ *bool = pflag.BoolP("help", "h", false, "")
	pflag.Parse()

	viper.SetDefault("PluginDir", "$XDG_CONFIG_HOME/tasker/autoload")
	viper.SetEnvPrefix("tasker")
	viper.AutomaticEnv()

	if *config != "" {
		f, err := os.Open(*config)
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

	if !(*noPlugins) {
		err := plugins.LoadPlugins()
		if err != nil {
			panic(err.Error())
		}
		defer plugins.UnloadPlugins()

		plugs := plugins.GetPlugins()
		for _, plug := range plugs {
			if commandsPlug, ok := plug.(plugins.TaskerCommand); ok {
				cmd.TaskerCmd.AddCommand(commandsPlug.Commands()...)
			}
		}
	}

	if err := cmd.Execute(); err != nil {
		log.Fatal(err.Error())
	}
}
