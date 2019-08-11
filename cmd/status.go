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
	"os"
	"sort"

	"github.com/oatmealraisin/tasker/pkg/models"
	"github.com/spf13/cobra"
)

// statusFlags isolates the flags specific to `tasker status`
var statusFlags struct {
	showFinished bool
	numShow      int
	tags         []string
}

// statusCmd is the cobra command for `tasker status`
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get the current to do list",
	Long:  `Get the current to do list`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.RunE(cmd, args); err != nil {
			log.Fatal(err.Error())
		}
	},
	RunE: status,
}

// Add `tasker status` to the command list, add `tasker status` flags
func init() {
	TaskerCmd.AddCommand(statusCmd)

	statusCmd.Flags().BoolVarP(&statusFlags.showFinished, "finished", "f", false, "Display even finished tasks.")
	statusCmd.Flags().IntVarP(&statusFlags.numShow, "number", "n", 10, "Number of tasks to display.")
	statusCmd.Flags().StringSliceVarP(&statusFlags.tags, "tag", "t", []string{}, "Give the status of tag or multiple tags.")
}

// status is the main function for the `tasker status` command. First we get the
// tasklist we're using for context (All by default, but could be within a list
// of tags). Then, we filter by IsNotFinished, IsNotRemoved. Finally, we sort
// by their score and print the first 10.
func status(cmd *cobra.Command, args []string) error {
	var err error
	if err = statusValidate(cmd, args); err != nil {
		return err
	}

	tasks := []uint64{}

	if len(statusFlags.tags) > 0 {
		tasks = db.GetByTags(statusFlags.tags)
	} else {
		tasks = db.GetAllTasks()
	}

	if len(tasks) == 0 {
		fmt.Printf("It doesn't look like you have anything to do!\n")
		return nil
	}

	filterList := models.FilterList{
		models.IsNotFinishedFilter(),
		models.IsNotRemovedFilter(),
		models.IsActiveFilter(),
		models.SizeIsNot(0),
		models.NoUnfinishedPrereqs(),
	}

	tasks = filterList.Apply(tasks, db.GetTask)

	sort.Slice(tasks, func(i, j int) bool {
		a, err := db.GetTask(tasks[i])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error sorting task %d: %s\n", tasks[i], err.Error())
			return false
		}

		b, err := db.GetTask(tasks[j])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error sorting task %d: %s\n", tasks[j], err.Error())
			return true
		}

		return a.Score() > b.Score()
	})

	models.PrintTasks(tasks[:statusFlags.numShow], db.GetTask)

	return nil
}

// statusValidate checks the flags and arguments for `tasker status` for errors
// or contradictions. It fills out anything left out.
func statusValidate(cmd *cobra.Command, args []string) error {
	// TODO: Implement
	return nil
}
