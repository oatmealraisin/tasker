package cmd

import (
	"log"
	"strconv"

	"github.com/oatmealraisin/tasker/pkg/models"
	"github.com/spf13/cobra"
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

func get(cmd *cobra.Command, args []string) error {
	var tasks []uint64
	if args[0] == "all" {
		tasks = db.GetAllTasks()
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

	models.PrintTasks(tasks, db.GetTask)

	return nil
}

func init() {
	TaskerCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
