package cmd

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/oatmealraisin/tasker/pkg/models"
	"github.com/spf13/cobra"
)

var (
	alsoChildren  bool
	alsoParents   bool
	alsoRelatives bool
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Retrieve specific information about specific tasks",
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
	RunE: get,
}

func init() {
	TaskerCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	getCmd.Flags().BoolVarP(&alsoChildren, "children", "C", false, "Print the children of the task as well.")
	getCmd.Flags().BoolVarP(&alsoParents, "parent", "p", false, "Print the parents of the task as well.")
}

func get(cmd *cobra.Command, args []string) error {
	var tasks []uint64
	if args[0] == "all" {
		allTasks := db.GetAllTasks()
		tasks = allTasks[:0]

		for _, u := range allTasks {
			if task, err := db.GetTask(u); err == nil {
				if !task.Removed && task.Added.Seconds < time.Now().Unix() {
					tasks = append(tasks, u)
				}
			} else {
				fmt.Fprintf(os.Stderr, err.Error())
			}
		}

	} else {
		uuid, err := strconv.ParseUint(args[0], 10, 64)
		if err != nil {
			return err
		}

		if _, err := db.GetTask(uuid); err != nil {
			return err
		}

		tasks = []uint64{uuid}
	}

	if alsoChildren {
		var children []uint64

		for _, uuid := range tasks {
			c, err := models.GetAllChildren(uuid, db.GetTask)
			if err != nil {
				return err
			}

			children = append(children, c...)
		}

		tasks = children
	}

	models.PrintTasks(tasks, db.GetTask)

	return nil
}
