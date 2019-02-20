package models

import (
	"fmt"
	"os"
)

type Filter func(task Task, get func(uuid uint64) (Task, error)) bool

func IsFinishedFilter(task Task, get func(uuid uint64) (Task, error)) bool {
	return task.Finished != nil
}

func IsNotFinishedFilter(task Task, get func(uuid uint64) (Task, error)) bool {
	return task.Finished == nil
}

func IsRemovedFilter(task Task, get func(uuid uint64) (Task, error)) bool {
	return task.Removed
}

func IsNotRemovedFilter(task Task, get func(uuid uint64) (Task, error)) bool {
	return !task.Removed
}

func HasDueDateFilter(task Task, get func(uuid uint64) (Task, error)) bool {
	return task.Due != nil
}

func DoesNotHaveDueDateFilter(task Task, get func(uuid uint64) (Task, error)) bool {
	return task.Due == nil
}

func NoUnfinishedPrereqs(task Task, get func(uuid uint64) (Task, error)) bool {
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

func UnfinishedPrereqs(task Task, get func(uuid uint64) (Task, error)) bool {
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
			if !filter(task, get) {
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
