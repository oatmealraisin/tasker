// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"log"

	"github.com/oatmealraisin/tasker/pkg/models"
	"github.com/spf13/cobra"
)

var (
	name string
	size string
	tags string
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "A brief description of your command",
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
	RunE: add,
}

func init() {
	TaskerCmd.AddCommand(addCmd)

	startCmd.Flags().StringVarP(&name, "name", "n", "", "Display name of the task")
	startCmd.Flags().StringVarP(&size, "size", "s", "", "Sizing for this task")
	startCmd.Flags().StringVarP(&tags, "tags", "t", "", "Tags to put this task in")
}

func add(cmd *cobra.Command, args []string) error {
	if err := validate(cmd, args); err != nil {
		return err
	}

	name := cmd.Flag("name").Value.String()
	size := cmd.Flag("size").Value.String()
	tags := cmd.Flag("tags").Value.String()

	newTask := &models.Task{
		Name: name,
		Size: size,
		Tags: tags.split(" "),
	}

	return nil
}

func validate(cmd *cobra.Command, args []string) error {
	// TODO: Implement
	if len(cmd.Flag("name").Value.String()) == 0 {
		return fmt.Errorf("Need to provide name to add new task")
	}
	return nil
}
