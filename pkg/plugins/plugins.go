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
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"plugin"

	"github.com/spf13/viper"
)

var loaded []TaskerPlugin

func loadPlugin(path string) (TaskerPlugin, error) {
	so, err := plugin.Open(path)
	if err != nil {
		return nil, fmt.Errorf("%s: %s\n", path, err.Error())
	}

	p, err := so.Lookup("Plugin")
	if err != nil {
		return nil, NewNoPluginError(path)
	}

	f, ok := p.(func() TaskerPlugin)
	if !ok {
		return nil, NewNotTaskerPluginError(path)
	}

	return f(), nil
}

func LoadPlugins() error {
	if len(loaded) != 0 {
		return fmt.Errorf("Plugins already loaded\n")
	}

	pluginDir := os.ExpandEnv(viper.GetString("PluginDir"))
	if pluginDir == "" {
		return fmt.Errorf("No plugin directory")
	}

	fns, err := ioutil.ReadDir(pluginDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Could not open plugin directory %s: %s\n", pluginDir, err.Error())
	}

	for _, fn := range fns {
		newPlugin, err := loadPlugin(filepath.Join(pluginDir, fn.Name()))
		if err != nil {
			return err
		}

		newPlugin.Initialize()

		loaded = append(loaded, newPlugin)
	}

	return nil
}

func UnloadPlugins() {
	for _, p := range loaded {
		p.Destroy()
	}
}

func GetPlugins() []TaskerPlugin {
	return loaded
}
