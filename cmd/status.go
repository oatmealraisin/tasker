package cmd

import (
	"fmt"
	"log"
	"math"
	"sort"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/oatmealraisin/tasker/pkg/models"
	"github.com/spf13/cobra"
)

var (
	showFinished bool
	numShow      int
)

// addCmd represents the add command
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

func init() {
	TaskerCmd.AddCommand(statusCmd)

	statusCmd.Flags().BoolVarP(&showFinished, "finished", "f", false, "Display even finished tasks.")
	statusCmd.Flags().IntVarP(&numShow, "number", "n", 10, "Number of tasks to display.")
}

func status(cmd *cobra.Command, args []string) error {
	if err := status_validate(cmd, args); err != nil {
		return err
	}

	// TODO: Find a way to go through only some of the tasks.. not all
	tasks := db.GetAllTasks()
	if len(tasks) == 0 {
		fmt.Println("It doesn't look like you have anything to do!")
		return nil
	}

	sort.Slice(tasks, func(i, j int) bool {
		a, err := db.GetTask(tasks[i])
		if err != nil {
			return false
		}

		b, err := db.GetTask(tasks[j])
		if err != nil {
			return true
		}

		return calc_task_score(a) > calc_task_score(b)
	})

	var selected []uint64

	i := 0
	for _, uuid := range tasks {
		if i >= numShow {
			break
		}

		task, err := db.GetTask(uuid)
		if err != nil {
			// TODO: Log
			fmt.Println("Could get get Task %d\n", uuid)
			continue
		}

		add_time, err := ptypes.Timestamp(task.Added)

		// Don't show tasks that are added in the future
		if err == nil && add_time.After(time.Now()) {
			y, m, d := add_time.Date()
			y_n, m_n, d_n := time.Now().Date()

			if !(y <= y_n && m <= m_n && d <= d_n) {
				continue
			}
		}

		// TODO: Implement subtasks
		// TODO: Implement dependencies
		if !task.Removed && task.Size != 0 {
			if len(task.Dependencies) != 0 {
				for _, j := range task.Dependencies {
					dep, err := db.GetTask(j)
					// Don't show a task if it has an unfinished dependency
					if err == nil && !dep.Removed {
						goto cont
					}
				}
			}

			// Confirm for selection, we can't fail after this
			selected = append(selected, task.Guid)
			i++
		}
	cont:
	}

	models.PrintTasks(selected, db.GetTask)
	return nil
}

func status_validate(cmd *cobra.Command, args []string) error {
	// TODO: Implement
	return nil
}

func calc_task_score(task models.Task) float64 {
	add_time, err := ptypes.Timestamp(task.Added)
	if err != nil {
		return 0.0
	}

	due_mod := 0.0
	if task.Due != nil {
		due_time, err := ptypes.Timestamp(task.Due)
		if err != nil {
			return 0.0
		}

		due_mod = 24.0 / math.Exp(time.Until(due_time).Hours())
	}

	add_mod := math.Max(1.0, math.Log(time.Since(add_time).Hours()/24.0))

	// TODO: Size and priority have a special relationship.. You want to do the
	// smallest, most important tasks first, followed by the hardest, most
	// important tasks. Change the formula to reflect this
	// TODO: Also, the age kind of changes the priority.. or at least makes the
	// priority less important. Should also change the formula to reflect this
	priority_mod := math.Pow(3.0, float64(task.Priority)) * 0.075
	size_mod := 0.5 * float64(task.Size)

	return due_mod + add_mod/(priority_mod+size_mod)
}
