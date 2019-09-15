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
	"log"

	"github.com/spf13/cobra"

	"github.com/oatmealraisin/tasker/pkg/plugins"
)

// pluginCmd represents the plugin command
var pluginCmd = &cobra.Command{
	Use:   "plugin",
	Short: "Add/remove plugins from tasker",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: runPlugin,
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install plugins to tasker.",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: install,
}

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall plugins from tasker.",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.RunE(cmd, args); err != nil {
			log.Fatal(err.Error())
		}
	},
	RunE: uninstall,
}

func init() {
	TaskerCmd.AddCommand(pluginCmd)
	pluginCmd.AddCommand(installCmd)
	pluginCmd.AddCommand(uninstallCmd)
}

// runPlugin does nothing, and is a placeholder for future functionality
func runPlugin(cmd *cobra.Command, args []string) error {
	return nil
}

// install parses the cobra command and installs a plugin.
func install(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("Need more arguments")
	} else if len(args) > 1 {
		return fmt.Errorf("Need less arguments")
	}

	return plugins.InstallFromFile(args[0])
}

// uninstall parses the cobra command and removes a plugin
func uninstall(cmd *cobra.Command, args []string) error {
	// TODO: Implement
	return nil
}
