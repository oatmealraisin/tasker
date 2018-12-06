package cmd

import (
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/golang/protobuf/ptypes"
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
		if tasks[i].Priority == tasks[j].Priority {
			i_time, err := ptypes.Timestamp(tasks[i].Added)
			if err != nil {
				return tasks[i].Priority < tasks[j].Priority
			}

			j_time, err := ptypes.Timestamp(tasks[j].Added)
			if err != nil {
				return tasks[i].Priority < tasks[j].Priority
			}

			i_t_mod := math.Max(1.0, math.Log(time.Since(i_time).Hours()/24.0))
			j_t_mod := math.Max(1.0, math.Log(time.Since(j_time).Hours()/24.0))

			return i_t_mod/float64(tasks[i].Size) > j_t_mod/float64(tasks[j].Size)
		}

		return tasks[i].Priority < tasks[j].Priority
	})

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
	fmt.Fprintln(w, "Name\tSize\tAge\tTags\t\t")

	i := 0
	for _, task := range tasks {
		if i >= numShow {
			break
		}

		add_time, err := ptypes.Timestamp(task.Added)

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
					if err == nil && !dep.Removed {
						goto cont
					}
				}
			}
			age := "?"
			if time_added, err := ptypes.Timestamp(task.Added); err == nil {
				age = strconv.Itoa(int(math.Floor(time.Since(time_added).Hours() / 24.0)))
			}

			i++

			var url string
			var nameString string
			if parent, err := db.GetTask(task.Parent); err == nil && parent.Guid != 0 {
				if parent.Url != "" {
					url = "(+)"
				}

				fmt.Fprintln(w, fmt.Sprintf("%s\t\t\t(%s)\t\t", parent.Name, url))
				nameString = fmt.Sprintf("└──> %s", task.Name)
				url = ""
			} else {
				nameString = task.Name
			}

			if task.Url != "" {
				url = "(+)"
			}

			fmt.Fprintln(w, fmt.Sprintf("%s\t%d\t%s days\t(%s)\t%s\t", nameString, task.Size, age, strings.Join(task.Tags, "|"), url))
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
