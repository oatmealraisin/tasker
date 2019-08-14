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
	"log"

	"github.com/spf13/cobra"
)

func (t *Today) run(cmd *cobra.Command, args []string) error {
	for !t.initialized {
	}

	err := t.update()
	if err != nil {
		return err
	}

	t.printAll()

	return nil
}

func (t *Today) add(cmd *cobra.Command, args []string) error {
	for !t.initialized {
	}

	return nil
}

func (t *Today) rm(cmd *cobra.Command, args []string) error {
	for !t.initialized {
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

	todayCmd.AddCommand(todayAddCmd)
	todayCmd.AddCommand(todayRmCmd)

	return []*cobra.Command{todayCmd}
}
