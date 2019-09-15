// Tasker - A pluggable task server for keeping track of all those To-Do's
// Today - A plugin for focusing on a subset of tasks just for today
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
package today

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/oatmealraisin/tasker/pkg/storage"
)

func (t *Today) Initialize() error {
	now := time.Now()
	yest := time.Now().Add(-time.Hour * 24)
	configDir := os.ExpandEnv("$XDG_CONFIG_HOME/tasker/today/")

	// Open our jsonFile
	jsonFile, err := os.Open(fmt.Sprintf("%s/data", configDir))
	if err != nil {
		jsonFile.Close()

		if !os.IsNotExist(err) {
			return err
		} else {
			t.Now = fmt.Sprintf("%d-%d-%d", now.Year(), now.Month(), now.Day())
			t.Yesterday = fmt.Sprintf("%d-%d-%d", yest.Year(), yest.Month(), yest.Day())
			t.Today = []uint64{}
			t.Tasks = make(map[string][]uint64)
			t.Initialized = true

			return nil
		}
	}
	defer jsonFile.Close()

	b, _ := ioutil.ReadAll(jsonFile)

	err = json.Unmarshal(b, &t.Tasks)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	t.Now = fmt.Sprintf("%d-%d-%d", now.Year(), now.Month(), now.Day())
	t.Yesterday = fmt.Sprintf("%d-%d-%d", yest.Year(), yest.Month(), yest.Day())
	t.Today = t.Tasks[t.Now]
	t.Initialized = true

	return nil
}

func (t *Today) Destroy() error {
	for !t.Initialized {
	}

	b, err := json.Marshal(t.Tasks)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	configDir := os.ExpandEnv("$XDG_CONFIG_HOME/tasker/today")
	err = ioutil.WriteFile(fmt.Sprintf("%s/data", configDir), b, 0644)
	if err != nil {
		fmt.Println(err.Error())
	}

	return nil
}

func (t *Today) Install() error {
	configDir := os.ExpandEnv("$XDG_CONFIG_HOME/tasker/today/")
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		return os.Mkdir(configDir, 0700)
	}

	return nil
}

func (t *Today) Uninstall() error {
	return nil
}

func (t *Today) Name() string {
	return "Today"
}

func (t *Today) Description() string {
	panic("not implemented")
}

func (t *Today) Help() string {
	panic("not implemented")
}

func (t *Today) Version() string {
	return "0.1.0"
}

func (t *Today) SetGetFunc(get storage.GetFunc) {
	t.Get = get
}
