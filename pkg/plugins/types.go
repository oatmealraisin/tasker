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
package plugins

import (
	"github.com/spf13/cobra"
)

/* TaskerPlugin is the base interface that all plugins have to implement. */
type TaskerPlugin interface {
	// Initialize is called whenever the plugin is loaded, normally at the
	// beginning of a command.
	Initialize() error
	/* Destroy is called when tasker stops running. Perform all cleanup needs. */
	Destroy() error
	/* Install is called once when the `tasker plugin install` command is run,
	after the plugin is built and autoloaded. */
	Install() error
	/* Uninstall is called once when the `tasker plugin uninstall` */
	Uninstall() error

	/* Name is used for formatting, you should return the pretty-name of your
	plugin. */
	Name() string
	/* Description is used for formatting, you should return a long-form
	description of what your plugin does. */
	Description() string
	/* Help is used for formatting, you should return configuration and command
	information about your plugin */
	Help() string
	/* Version is used for comparing installed plugins. */
	Version() string
}

type TaskerCommand interface {
	Commands() []*cobra.Command
}
