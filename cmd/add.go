package cmd

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/oatmealraisin/tasker/pkg/models"
	"github.com/oatmealraisin/tasker/pkg/storage"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	name       string
	size       int
	tags       string
	priority   int
	url        string
	importFile string
	dryRun     bool
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
	// TODO:
	// PreRun: validate,
}

func init() {
	TaskerCmd.AddCommand(addCmd)

	addCmd.Flags().StringVarP(&name, "name", "n", "", "Display name of the task")
	addCmd.Flags().IntVarP(&size, "size", "s", 0, "Sizing for this task")
	addCmd.Flags().StringVarP(&tags, "tags", "t", "", "Tags to put this task in")
	addCmd.Flags().IntVarP(&priority, "priority", "p", -1, "The priority of this task, how important it is.")
	addCmd.Flags().StringVarP(&url, "url", "u", "", "Any URL resource associated with this tasks, such as an article.")
	addCmd.Flags().StringVarP(&importFile, "from-file", "f", "", "Import tasks from a file. Can be csv.")
	addCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Go through the steps but do nothing.")
}

func add(cmd *cobra.Command, args []string) error {
	if err := validate(cmd, args); err != nil {
		return err
	}

	tasks := []models.Task{}

	if importFile != "" {
		f, err := ioutil.ReadFile(importFile)
		if err != nil {
			return err
		}
		input := string(f)

		records, err := csv.NewReader(strings.NewReader(input)).ReadAll()
		if err != nil {
			return err
		}

		for i, record := range records {
			if i == 0 {
				continue
			}

			newTask, err := storage.TaskFromCsv(record)
			if err != nil {
				return err
			}
			tasks = append(tasks, newTask)
		}

	} else {
		newTask := models.Task{
			Name: name,
			Size: uint32(size),
			Tags: strings.Split(tags, viper.GetString("Delim")),
		}

		tasks = append(tasks, newTask)
	}

	// TODO: Only submit if all tasks are valid
	for _, task := range tasks {
		//	fmt.Printf("%d, %s, %d\n", task.Guid, task.Name, task.Size)
		err := db.CreateTask(task)
		if err != nil {
			return err
		}
	}

	return nil
}

func validate(cmd *cobra.Command, args []string) error {
	// TODO: Implement
	if importFile != "" && (name != "" || size != 0 || tags != "" || priority < 0 || url != "") {
		//return fmt.Errorf("-f/--from-file cannot be used with other flags.")
	}

	if importFile != "" {
		return nil
	}

	// TODO: Check for import_file existence and readable

	if len(name) == 0 {
		return fmt.Errorf("Need to provide name to add new task")
	}

	if size == 0 {

	}
	return nil
}
