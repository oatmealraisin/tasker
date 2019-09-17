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
	"time"

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

// A utility function for standard today usage. Will either print the content of
// the tasks registered for the current day, the content of the tasks registered
// for the next previously registered day, or nothing if no previous date has
// been registered.
func (t *Today) printToday() {
	t.printDate(t.Now)

	if len(t.Today) == 0 {
		prev := time.Now()
		for i := 0; i < 14; i++ {
			prev = prev.Add(-time.Hour * 24)
			prevYMD := fmt.Sprintf(
				"%d-%d-%d",
				prev.Year(),
				prev.Month(),
				prev.Day(),
			)
			if prevTasks, ok := t.Tasks[prevYMD]; ok && len(prevTasks) > 0 {
				fmt.Printf("Here's what happened %s:\n", prevYMD)
				t.printDate(prevYMD)
				return
			}
		}
	}

	return
}

// A generic function for printing any date from the days registered with today.
// Will either print the content of the days tasks using tasker libraries, or
// will print out a message telling the user that there is nothing to do.
// The date must be formatted '2006-1-2'.
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
