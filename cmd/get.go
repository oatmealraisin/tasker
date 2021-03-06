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

// getFlags is a utility struct for isolating variables specific to the `get`
// command
var getFlags struct {
	alsoChildren  bool
	alsoParents   bool
	alsoRelatives bool

	getAll          bool
	uuid            uint64
	tags            []string
	tagsOpt         []string
	includeFinished bool
	includeRemoved  bool
	dueBefore       string
	createdAfter    string
	tCreatedAfter   time.Time
	url             bool
	name            bool
}

// getCmd represents the get command, mostly Cobra boilerplate
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

	getCmd.Flags().BoolVarP(&getFlags.alsoChildren, "children", "C", false, "Print the children of the task as well.")
	getCmd.Flags().BoolVarP(&getFlags.alsoParents, "parent", "p", false, "Print the parents of the task as well.")
	// TODO: Implement a fuzzy finder for names
	//getCmd.Flags().StringVarP(&name, "name", "n", "", "Get tasks with a similar name")
	getCmd.Flags().StringSliceVarP(&getFlags.tags, "tag", "t", []string{}, "Get tasks from a tag. Can be invoked more than once to specify multiple tags.")
	getCmd.Flags().StringSliceVar(&getFlags.tagsOpt, "has-tag", []string{}, "Get tasks from a tag. Can be invoked more than once to specify multiple tags.")
	getCmd.Flags().Uint64VarP(&getFlags.uuid, "uuid", "u", 0, "Get tasks with matching uuid. Will only return one task.")
	getCmd.Flags().BoolVar(&getFlags.includeFinished, "finished", false, "Also give tasks that have been finished.")
	getCmd.Flags().BoolVar(&getFlags.includeRemoved, "removed", false, "Also give tasks that have been removed.")
	getCmd.Flags().BoolVar(&getFlags.url, "url", false, "Print the URL associated with the task.")
	getCmd.Flags().BoolVar(&getFlags.name, "name", false, "Print the name of the task only.")
	getCmd.Flags().StringVar(&getFlags.dueBefore, "due-before", "", "Only show tasks due before a certain date.")
	getCmd.Flags().StringVar(&getFlags.createdAfter, "created-after", "", "Only show tasks created after a certain date.")
}

func get(cmd *cobra.Command, args []string) error {
	var tasks []uint64

	if err := validateGet(cmd, args); err != nil {
		return err
	}

	if getFlags.getAll {
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
	}

	if getFlags.uuid != 0 {
		_, err := db.GetTask(getFlags.uuid)
		if err != nil {
			return err
		}

		tasks = []uint64{getFlags.uuid}
	} else if cmd.Flag("tag").Changed {
		tasks = db.GetByTags(getFlags.tags)
	}

	if getFlags.alsoChildren {
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

	filterList := models.FilterList{}

	if !getFlags.includeFinished {
		filterList = append(filterList, models.IsNotFinishedFilter())
	}

	if !getFlags.includeRemoved {
		filterList = append(filterList, models.IsNotRemovedFilter())
	}

	if getFlags.createdAfter != "" {
		filterList = append(filterList, models.CreatedAfter(getFlags.tCreatedAfter))
	}

	tasks = filterList.Apply(tasks, db.GetTask)

	if getFlags.url {
		for _, uuid := range tasks {
			task, err := db.GetTask(uuid)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Could not get task %d.\n%s", uuid, err.Error())
				continue
			}

			if task.Url != "" {
				fmt.Printf("%s\n", task.Url)
			}
		}
	} else if getFlags.name {
		for _, uuid := range tasks {
			task, err := db.GetTask(uuid)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Could not get task %d.\n%s", uuid, err.Error())
				continue
			}

			fmt.Printf("%s\n", task.Name)
		}
	} else {
		models.PrintTasks(tasks, db.GetTask)
	}

	return nil
}

// validate checks to make sure there aren't any contradictions or out of bound
// fields in the invocation.
func validateGet(cmd *cobra.Command, args []string) error {
	var err error

	if cmd.Flag("uuid").Changed {
		if len(getFlags.createdAfter)+len(getFlags.tags)+len(args)+len(getFlags.tagsOpt) != 0 || cmd.Flag("include-finished").Changed {
			return fmt.Errorf("Cannot specify multiple filters if UUID is given.\n")
		}

		getFlags.includeFinished = true
	}

	if getFlags.createdAfter != "" {
		getFlags.tCreatedAfter, err = time.Parse("2006-01-02", getFlags.createdAfter)
		if err != nil {
			return err
		}
	}

	if len(args) == 1 {
		if args[0] == "all" {
			if len(getFlags.tags)+len(getFlags.tagsOpt) != 0 {
				return fmt.Errorf("All means all! Specify tag filters without 'all'\n")
			}

			getFlags.getAll = true
		} else if uuid, err := strconv.ParseUint(args[0], 10, 64); err == nil {
			getFlags.uuid = uuid
			getFlags.includeFinished = true
		}
	}

	return nil
}
