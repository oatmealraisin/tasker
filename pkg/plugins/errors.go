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

import "fmt"

type NoPluginError struct {
	Name string
}

func (e *NoPluginError) Error() string {
	return fmt.Sprintf("ERROR: 'Plugin %s' has not exported a Plugin", e.Name)
}

func NewNoPluginError(pluginName string) error {
	return &NoPluginError{Name: pluginName}
}

type NotTaskerPluginError struct {
	Name string
}

func (e *NotTaskerPluginError) Error() string {
	msg := "Plugin '%s' does not implement the TaskerPlugin interface!"
	return fmt.Sprintf(msg, e.Name)
}

func NewNotTaskerPluginError(pluginName string) error {
	return &NotTaskerPluginError{Name: pluginName}
}
