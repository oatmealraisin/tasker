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
package models

import (
	"fmt"
	"os"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/spf13/viper"
)

type Filter struct {
	Apply func(task Task, get func(uuid uint64) (Task, error)) bool

	opt bool
}

func (f Filter) NotOpt() Filter {
	f.opt = false

	return f
}

func (f Filter) Opt() Filter {
	f.opt = true

	return f
}

func IsFinishedFilter() (result Filter) {
	result.Apply = func(task Task, get func(uuid uint64) (Task, error)) bool {
		return task.Finished != nil
	}

	return result
}

func IsNotFinishedFilter() (result Filter) {
	result.Apply = func(task Task, get func(uuid uint64) (Task, error)) bool {
		return task.Finished == nil
	}

	return
}

func IsRemovedFilter() (result Filter) {
	result.Apply = func(task Task, get func(uuid uint64) (Task, error)) bool {
		return task.Removed
	}

	return result
}

func IsNotRemovedFilter() (result Filter) {
	result.Apply = func(task Task, get func(uuid uint64) (Task, error)) bool {
		return !task.Removed
	}

	return result
}

func HasDueDateFilter() (result Filter) {
	result.Apply = func(task Task, get func(uuid uint64) (Task, error)) bool {
		return task.Due != nil
	}

	return result
}

func DoesNotHaveDueDateFilter() (result Filter) {
	result.Apply = func(task Task, get func(uuid uint64) (Task, error)) bool {
		return task.Due == nil
	}

	return result
}

func IsStale() (result Filter) {
	result.Apply = func(task Task, get func(uuid uint64) (Task, error)) bool {
		add_time, err := ptypes.Timestamp(task.Added)
		if err != nil {
			return false
		}

		return time.Since(add_time).Hours() >= float64(viper.GetInt("StaleTime"))
	}

	return result
}

func IsNotStale() (result Filter) {
	result.Apply = func(task Task, get func(uuid uint64) (Task, error)) bool {
		return !IsStale().Apply(task, get)
	}

	return result
}

func NoUnfinishedPrereqs() (result Filter) {
	result.Apply = func(task Task, get func(uuid uint64) (Task, error)) bool {
		for _, uuid := range task.Dependencies {
			prereq, err := get(uuid)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Could not get task %d\n", uuid)
				return false
			}

			if prereq.Finished == nil {
				return false
			}
		}

		return true
	}

	return result
}

func UnfinishedPrereqs() (result Filter) {
	result.Apply = func(task Task, get func(uuid uint64) (Task, error)) bool {
		for _, uuid := range task.Dependencies {
			prereq, err := get(uuid)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Could not get task %d\n", uuid)
				return false
			}

			if prereq.Finished == nil {
				return true
			}
		}

		return false
	}

	return result
}

func IsActiveFilter() (result Filter) {
	result.Apply = func(task Task, get func(uuid uint64) (Task, error)) bool {
		add_time, err := ptypes.Timestamp(task.Added)
		return err == nil && !add_time.After(time.Now())
	}

	return result
}

func IsNotActiveFilter() (result Filter) {
	result.Apply = func(task Task, get func(uuid uint64) (Task, error)) bool {
		add_time, err := ptypes.Timestamp(task.Added)
		return err != nil || add_time.After(time.Now())
	}

	return result
}

func SizeIsExactly(size uint32) (result Filter) {
	result.Apply = func(task Task, get func(uuid uint64) (Task, error)) bool {
		return task.Size == size
	}

	return result
}

func SizeIsNot(size uint32) (result Filter) {
	result.Apply = func(task Task, get func(uuid uint64) (Task, error)) bool {
		return task.Size != size
	}

	return result
}

func HasParent() (result Filter) {
	result.Apply = func(task Task, get func(uuid uint64) (Task, error)) bool {
		return task.Parent != 0
	}

	return result
}

func HasChildren() (result Filter) {
	result.Apply = func(task Task, get func(uuid uint64) (Task, error)) bool {
		return len(task.Subtasks) != 0
	}

	return result
}

func CreatedAfter(date time.Time) (result Filter) {
	result.Apply = func(task Task, get func(uuid uint64) (Task, error)) bool {
		added, err := ptypes.Timestamp(task.Added)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting add date of Task %d (%s): %s", task.Guid, task.Name, err.Error())
		}

		return added.After(date)
	}

	return result
}

type FilterList []Filter

// Apply returns a subset of `uuids` that pass through all the filters in the
// list
func (f FilterList) Apply(uuids []uint64, get func(uuid uint64) (Task, error)) []uint64 {
	result := []uint64{}
	for _, uuid := range uuids {
		task, err := get(uuid)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not get task %d\n", uuid)
			goto reject
		}

		for _, filter := range f {
			if !filter.Apply(task, get) {
				goto reject
			}
		}

		result = append(result, uuid)

	reject:
	}

	return result
}

func (f FilterList) parallelApply(uuid []uint64, get func(uuid uint64) (Task, error)) []uint64 {
	return []uint64{}
}
