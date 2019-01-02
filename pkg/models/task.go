package models

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/golang/protobuf/ptypes"
	"golang.org/x/crypto/ssh/terminal"
)

func (task Task) Score() float64 {
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

type Tasks []Task

func (tasks Tasks) Pretty() string {
	return ""
}

func PrintTasks(tasks []uint64, get func(uuid uint64) (Task, error)) {
	var nameString string
	var url string

	termWidth, _, err := terminal.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		panic(err.Error())
	}

	// TODO: Modify to start removing columns to fit screen
	maxSize := termWidth - 55

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
	fmt.Fprintln(w, "Name\tSize\tAge\tDue\tTags\t\t")

	for _, uuid := range tasks {
		task, err := get(uuid)
		if err != nil {
			// TODO: log
			fmt.Println("Error getting task %d\n", uuid)
			continue
		}

		if task.Parent != 0 {
			parent, err := get(task.Parent)
			if err != nil {
				// TODO: Log
				fmt.Println("Error getting parent of %d: %d\n", task.Guid, task.Parent)
				return
			}

			nameString = fmt.Sprintf("└──> %s", task.Name)

			var parentString string
			if parent.Parent != 0 {
				grand, err := get(parent.Parent)
				if err != nil {
					// TODO: Log
					fmt.Println("Error getting grandparent of %d: %d\n", parent.Parent, task.Guid)
					return
				}

				if grand.Url != "" {
					url = "(+)"
				}

				if grand.Parent != 0 {
					for grand.Parent != 0 {
						new_grand, err := get(grand.Parent)
						if err != nil {
							// TODO: Log
							fmt.Println("Error getting grandparent of %d: %d\n", task.Guid, grand.Parent)
							return
						}

						grand = new_grand
					}
				}

				grandName := grand.Name
				if len(grandName) > maxSize {
					grandName = fmt.Sprintf("%s ...", grandName[:maxSize])
				}

				fmt.Fprintln(w, fmt.Sprintf("%s\t\t\t\t%s\t", grandName, url))
				parentString = fmt.Sprintf("└┬─> %s", parent.Name)
				nameString = fmt.Sprintf(" %s", nameString)
			} else {
				parentString = parent.Name
			}

			if parent.Url != "" {
				url = "(+)"
			}

			if len(parentString) > maxSize {
				parentString = fmt.Sprintf("%s ...", parentString[:maxSize])
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

		if len(nameString) > maxSize {
			nameString = fmt.Sprintf("%s ...", nameString[:maxSize])
		}

		fmt.Fprintln(w, fmt.Sprintf("%s\t %d %d\t%s\t%s\t(%s)\t%s\t", nameString, task.Size, task.Priority, age, due, strings.Join(task.Tags, "|"), url))
	}

	w.Flush()

	return
}
