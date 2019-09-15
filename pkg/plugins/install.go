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
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/spf13/viper"
)

// InstallFromFile installs a plugin to the PluginDir, and runs the Install
// function of the plugin
func InstallFromFile(fn string) error {
	pluginDir := viper.GetString("PluginDir")
	if pluginDir == "" {
		return fmt.Errorf("No PluginDir set.")
	}

	pluginDir = os.ExpandEnv(pluginDir)

	f, err := os.Open(fn)
	if err != nil {
		return fmt.Errorf("Could not open given file: %s", err.Error())
	}

	newPlugin, err := loadPlugin(f.Name())
	if err != nil {
		return err
	}

	err = newPlugin.Install()
	if err != nil {
		return err
	}

	dst, err := os.Create(filepath.Join(pluginDir, path.Base(f.Name())))
	if err != nil {
		return fmt.Errorf("Could not create file in autoload directory: %s", err.Error())
	}
	defer dst.Close()

	_, err = io.Copy(dst, f)
	if err != nil {
		return err
	}

	err = dst.Close()
	if err != nil {
		return err
	}

	return nil
}

func InstallFromGit(url string) (*TaskerPlugin, error) {
	// TODO: Check for .so in tld of project
	// TODO: Check for Makefile target in tld of project
	return nil, nil
}

func InstallFromUrl(url string) (*TaskerPlugin, error) {
	// TODO: check for git

	return nil, nil
}
