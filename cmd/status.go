package cmd

import (
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
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

	tasks, err := db.GetAllTasks()
	if err != nil {
		return err
	}

	sort.Slice(tasks, func(i, j int) bool {
		return calc_task_score(tasks[i]) > calc_task_score(tasks[j])
	})

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
	fmt.Fprintln(w, "Name\tSize\tAge\tDue\tTags\t\t")

	i := 0
	for _, task := range tasks {
		if i >= numShow {
			break
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
			i++
			printTask(w, task)
		}
	cont:
	}

	w.Flush()

	return nil
}

func status_validate(cmd *cobra.Command, args []string) error {
	// TODO: Implement
	return nil
}

func printTask(w io.Writer, task models.Task) error {
	var nameString string

	if task.Parent != 0 {
		parent, err := db.GetTask(task.Parent)
		if err != nil {
			return err
		}

		nameString = fmt.Sprintf("└──> %s", task.Name)

		var parentString string
		if parent.Parent != 0 {
			grand, err := db.GetTask(parent.Parent)
			if err != nil {
				return err
			}

			if grand.Url != "" {
				url = "(+)"
			}

			if grand.Parent != 0 {
				for grand.Parent != 0 {
					new_grand, err := db.GetTask(grand.Parent)
					if err != nil {
						return err
					}

					grand = new_grand
				}
			}

			fmt.Fprintln(w, fmt.Sprintf("%s\t\t\t\t%s\t", grand.Name, url))
			parentString = fmt.Sprintf("└┬─> %s", parent.Name)
			nameString = fmt.Sprintf(" %s", nameString)
		} else {
			parentString = parent.Name
		}

		if parent.Url != "" {
			url = "(+)"
		}

		fmt.Fprintln(w, fmt.Sprintf("%s\t\t\t\t%s\t", parentString, url))
		url = ""
	} else {
		nameString = task.Name
	}

	var due string
	if task.Due != nil {
		if due_time, err := ptypes.Timestamp(task.Due); err == nil {
			num_days := int(math.Floor(time.Until(due_time).Hours() / 24.0))
			if num_days == 0 {
				due = "Today"
			} else if num_days == 1 {
				due = "Tmrw"
			} else {
				due = fmt.Sprintf("%sd", strconv.Itoa(num_days))
			}
		}
	}

	age := "?"
	if time_added, err := ptypes.Timestamp(task.Added); err == nil {
		num_days := int(math.Floor(time.Since(time_added).Hours() / 24.0))
		if num_days == 0 {
			age = "New"
		} else if num_days < 30 {
			age = fmt.Sprintf("%sd", strconv.Itoa(num_days))
		} else if num_days < 360 {
			age = fmt.Sprintf("%sm", strconv.Itoa(num_days/30.0))
		} else {
			age = ">1y"
		}
	}

	var url string
	if task.Url != "" {
		url = "(+)"
	}

	fmt.Fprintln(w, fmt.Sprintf("%s\t %d %d\t%s\t%s\t(%s)\t%s\t", nameString, task.Size, task.Priority, age, due, strings.Join(task.Tags, "|"), url))

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
