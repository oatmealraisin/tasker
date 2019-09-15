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
	"fmt"

	"github.com/oatmealraisin/tasker/pkg/models"
	"github.com/oatmealraisin/tasker/pkg/storage"
)

type Today struct {
	Tasks       map[string][]uint64
	Today       []uint64 `json: "-"`
	Now         string   `json: "-"`
	Yesterday   string   `json: "-"`
	Initialized bool     `json: "-"`

	Get storage.GetFunc `json: "-"`
}

func (t *Today) printToday() {

	t.printDate(t.Now)

	if len(t.Today) == 0 {
		if yest, ok := t.Tasks[t.Yesterday]; ok && len(yest) > 0 {
			fmt.Println("Here's what happened yesterday:")
			t.printDate(t.Yesterday)
		}
	}

	return
}

func (t *Today) printDate(day string) {
	// TODO: Check validity of day
	if uuids, ok := t.Tasks[day]; ok {
		if t.Get == nil {
			fmt.Println(uuids)
		} else {
			models.PrintTasks(uuids, t.Get)
		}
	} else {
		fmt.Printf("That's all for %s\n", day)
	}
}
