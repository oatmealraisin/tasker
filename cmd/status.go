package cmd

import (
	"fmt"
	"log"
	"os"
	"sort"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/oatmealraisin/tasker/pkg/models"
	"github.com/spf13/cobra"
)

var statusFlags struct {
	showFinished bool
	numShow      int
	tags         []string
}

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

	statusCmd.Flags().BoolVarP(&statusFlags.showFinished, "finished", "f", false, "Display even finished tasks.")
	statusCmd.Flags().IntVarP(&statusFlags.numShow, "number", "n", 10, "Number of tasks to display.")
	statusCmd.Flags().StringSliceVarP(&statusFlags.tags, "tag", "t", []string{}, "Give the status of tag or multiple tags.")
}

func status(cmd *cobra.Command, args []string) error {
	if err := statusValidate(cmd, args); err != nil {
		return err
	}

	tasks := []uint64{}

	if len(statusFlags.tags) > 0 {
		tasks = db.GetByTags(statusFlags.tags)
	} else {
		tasks = db.GetAllTasks()
	}

	if len(tasks) == 0 {
		fmt.Println("It doesn't look like you have anything to do!")
		return nil
	}

	sort.Slice(tasks, func(i, j int) bool {
		a, err := db.GetTask(tasks[i])
		if err != nil {
			fmt.Println("Error: Couldn't get Task %d: %s", tasks[i], err.Error())
			return false
		}

		b, err := db.GetTask(tasks[j])
		if err != nil {
			fmt.Println("Error: Couldn't get Task %d: %s", tasks[j], err.Error())
			return true
		}

		return a.Score() > b.Score()
	})

	var selected []uint64

	i := 0
	for _, uuid := range tasks {
		if i >= statusFlags.numShow {
			break
		}

		task, err := db.GetTask(uuid)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could get get Task %d\n", uuid)
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

func statusValidate(cmd *cobra.Command, args []string) error {
	// TODO: Implement
	return nil
}
