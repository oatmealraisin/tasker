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
	"log"
	"sort"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

func (t *Today) run(cmd *cobra.Command, args []string) error {
	for !t.Initialized {
	}

	if len(args) > 0 {
		_, err := time.Parse("2006-1-2", args[0])
		if err != nil {
			return err
		}
		t.printDate(args[0])
	} else {
		t.printToday()
	}

	return nil
}

func (t *Today) add(cmd *cobra.Command, args []string) error {
	for !t.Initialized {
	}

	add_date := t.Now
	for _, x := range args {
		uuid, err := strconv.ParseUint(x, 10, 64)
		if err != nil {
			_, err = time.Parse("2006-1-2", x)
			if err != nil {
				return err
			}

			add_date = x
			continue
		}

		if _, ok := t.Tasks[add_date]; ok {
			t.Tasks[add_date] = append(t.Tasks[add_date], uuid)
		} else {
			t.Tasks[add_date] = []uint64{uuid}
		}
	}

	return nil
}

func (t *Today) rm(cmd *cobra.Command, args []string) error {
	for !t.Initialized {
	}

	return nil
}

func (t *Today) list(cmd *cobra.Command, args []string) error {
	keys := []string{}
	for k := range t.Tasks {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(a, b int) bool {
		a_time, err := time.Parse("2006-1-2", keys[a])
		if err != nil {
			return false
		}

		b_time, err := time.Parse("2006-1-2", keys[b])
		if err != nil {
			return true
		}

		return b_time.Before(a_time)
	})

	for i := 0; i < 10; i++ {
		fmt.Printf("%s: %d\n", keys[i], len(t.Tasks[keys[i]]))
	}

	return nil
}

func (t *Today) Commands() []*cobra.Command {
	todayCmd := &cobra.Command{
		Use:   "today",
		Short: "A plugin for focusing on a subset of tasks just for today",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := cmd.RunE(cmd, args); err != nil {
				log.Fatal(err.Error())
			}
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return t.run(cmd, args)
		},
	}

	todayAddCmd := &cobra.Command{
		Use:   "add",
		Short: "A plugin for focusing on a subset of tasks just for today",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := cmd.RunE(cmd, args); err != nil {
				log.Fatal(err.Error())
			}
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return t.add(cmd, args)
		},
	}

	todayRmCmd := &cobra.Command{
		Use:   "rm",
		Short: "A plugin for focusing on a subset of tasks just for today",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := cmd.RunE(cmd, args); err != nil {
				log.Fatal(err.Error())
			}
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return t.rm(cmd, args)
		},
	}

	todayListCmd := &cobra.Command{
		Use:   "list",
		Short: "A plugin for focusing on a subset of tasks just for today",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := cmd.RunE(cmd, args); err != nil {
				log.Fatal(err.Error())
			}
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return t.list(cmd, args)
		},
	}

	todayCmd.AddCommand(todayAddCmd)
	todayCmd.AddCommand(todayRmCmd)
	todayCmd.AddCommand(todayListCmd)

	return []*cobra.Command{todayCmd}
}
