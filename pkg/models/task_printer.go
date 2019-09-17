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
	"io"
	math "math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	minNameLenShort int = 50
	minSizeLen      int = minNameLenShort + 15
	minDueLenShort  int = minSizeLen + 14
	minAgeLen       int = minDueLenShort + 7
	minURLLen       int = minAgeLen + 7
	minTagLen       int = minURLLen + 50
	minDueLenLong   int = minTagLen + 10
	minNameLenLong  int = minDueLenLong + 10
	minTagLenLong   int = minNameLenLong + 20

	longNameLen int = 50
	longTagLen  int = 20

	removeTag *regexp.Regexp = regexp.MustCompile("[^|]*(\\|\\.\\.\\.)?$")
)

// getTermWidth is a utility function for understanding where we will be forced
// to newline
func getTermWidth() int {
	result, _, err := terminal.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return 200
	}

	return result
}

// PrintTasks is shorthand for creating a TaskTree from a list of UUIDs, then
// printing the tree.
func PrintTasks(tasks []uint64, get func(uuid uint64) (Task, error)) {
	if len(tasks) == 0 {
		return
	}

	tt := CreateTaskTree(tasks, get)
	if tt == nil {
		fmt.Fprintf(os.Stderr, "Error creating task tree\n")
		return
	}

	tt.Print(get)
}

// stringify preps a map to be passed to printStringyTask. This intermediate
// step is useful for editing the fields before the final print.
//
// TODO: This is a pretty wonky paradigm. We would do well to make it more
// intuitive
func (task Task) stringify() map[string]string {
	result := make(map[string]string)

	result["due"] = ""
	if task.Due != nil {
		if due_time, err := ptypes.Timestamp(task.Due); err == nil {
			due_time = due_time.In(time.Now().Location())
			num_days := time.Until(due_time).Hours() / 24.0
			if num_days < 0.0 {
				result["due"] = fmt.Sprintf("%04d-%02d-%02d %02d:%02d", due_time.Year(), due_time.Month(), due_time.Day(), due_time.Hour(), due_time.Minute())
			} else if num_days < 1.0 {
				result["due"] = fmt.Sprintf("%dh", int(time.Until(due_time).Hours()+0.5))
			} else if num_days < 2.0 {
				result["due"] = "Tmrw"
			} else if num_days < 30.0 {
				result["due"] = fmt.Sprintf("%dd", int(num_days))
			} else if num_days < 360.0 {
				result["due"] = fmt.Sprintf("%dm", int(num_days/30))
			} else {
				result["due"] = ">1y"
			}
		}
	}

	result["age"] = "?"
	if time_added, err := ptypes.Timestamp(task.Added); err == nil {
		num_days := int(math.Floor(time.Since(time_added).Hours() / 24.0))
		if num_days == 0 {
			result["age"] = "New"
		} else if num_days < 30 {
			result["age"] = fmt.Sprintf("%dd", num_days)
		} else if num_days < 360 {
			result["age"] = fmt.Sprintf("%dm", num_days/30)
		} else {
			result["age"] = ">1y"
		}
	}

	result["name"] = task.Name
	if task.Finished != nil {
		result["finished"] = "?"
		result["name"] = fmt.Sprintf("\x1B[9m%s\x1B[0m", result["name"])

		if time_finished, err := ptypes.Timestamp(task.Finished); err == nil {
			num_days := int(math.Floor(time.Since(time_finished).Hours() / 24.0))
			if num_days == 0 {
				result["finished"] = "Today"
			} else if num_days < 30 {
				result["finished"] = fmt.Sprintf("%dd ago", num_days)
			} else if num_days < 360 {
				result["finished"] = fmt.Sprintf("%dm ago", num_days/30)
			} else {
				result["finished"] = ">1y ago"
			}
		}
	}

	result["url"] = ""
	if task.Url != "" {
		result["url"] = "(+)"
	}

	if viper.GetBool("debug") {
		result["name"] = fmt.Sprintf("%d %s", task.Guid, result["name"])
	}

	result["tags"] = strings.Join(task.Tags, "|")
	result["size"] = strconv.FormatUint(uint64(task.Size), 10)
	result["priority"] = strconv.FormatUint(uint64(task.Priority), 10)

	return result
}

// TODO: Merge this with printStringyTask somehow.
// printColumns Prints the columns that correlate to the printStringyTask
// function. It uses the same calculations for column removal
func printColumns(w io.Writer) {
	defer fmt.Fprintln(w, "")
	p := func(s string) {
		fmt.Fprintf(w, fmt.Sprintf("%s\t", s))
	}

	termWidth := getTermWidth()

	if termWidth > len("Name") {
		p("Name")
	}

	if termWidth > minSizeLen {
		p("Size")
	}

	if termWidth > minAgeLen {
		p("Age")
	}

	if termWidth > minDueLenShort {
		p("Due")
	}

	if termWidth > minTagLen {
		p("Tag")
	}

	if termWidth > minURLLen {
		p("")
	}
}

func printStringyTask(w io.Writer, sfyTask map[string]string) {
	result := ""
	p := func(s string) {
		result = fmt.Sprintf("%s%s\t", result, s)
	}

	termWidth := getTermWidth()
	if termWidth < 10 {
		fmt.Fprintln(w, "...")
		return
	}

	name := sfyTask["name"]

	// If the name was crossed out (if it was finished), we still want it to
	// be aligned with task names that are not finished
	hasEscape := strings.Contains(name, "[0m")

	nameCutoff := 0
	if termWidth > minNameLenLong {
		nameCutoff = int(math.Max(float64(minNameLenShort), float64(termWidth-55)))
	} else if termWidth > minNameLenShort {
		nameCutoff = longNameLen
	} else {
		nameCutoff = termWidth - 4
	}

	if hasEscape {
		name = name[:len(name)-4]
		nameCutoff = nameCutoff + 4
	}

	if len(name) > nameCutoff {
		bName := []byte(fmt.Sprintf("%s", name[:nameCutoff]))
		bName[nameCutoff-1] = '.'
		bName[nameCutoff-2] = '.'
		bName[nameCutoff-3] = '.'
		bName[nameCutoff-4] = ' '
		name = string(bName)
	}

	name = fmt.Sprintf("%s\x1B[0m", name)

	p(name)

	if termWidth > minSizeLen {
		p(fmt.Sprintf(" %s", sfyTask["size"]))
	}

	if termWidth > minAgeLen {
		p(sfyTask["age"])
	}

	if termWidth > minDueLenLong {
		p(sfyTask["due"])
	} else if termWidth > minDueLenShort {
		if len(sfyTask["due"]) > 3 {
			p(sfyTask["due"][:len(sfyTask["due"])-6])
		} else {
			p("")
		}
	}

	if termWidth > minTagLen {
		tags := sfyTask["tags"]

		if termWidth < minTagLenLong && len(tags) > longTagLen {
			maxTagLen := termWidth - minTagLen
			for len(tags) > longTagLen {
				if removeTag.MatchString(tags) {
					tags = removeTag.ReplaceAllString(tags, "...")
				} else {
					tags = fmt.Sprintf("%s ...", tags[:maxTagLen])
				}
			}
		}

		p(fmt.Sprintf("(%s)", tags))
	}

	if termWidth > minURLLen {
		p(sfyTask["url"])
	}

	fmt.Fprintln(w, result)
}
