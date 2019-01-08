package cmd

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	"github.com/oatmealraisin/tasker/pkg/models"
	"github.com/oatmealraisin/tasker/pkg/storage"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	name  string
	size  int
	tags  string
	due_s string
	due   time.Time
	// TODO: Fix this
	due_p      *timestamp.Timestamp
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
	addCmd.Flags().StringVarP(&due_s, "due", "d", "", "Due date of this task, in YYYY-MM-DD form.")
	addCmd.Flags().IntVarP(&priority, "priority", "p", -1, "The priority of this task, how important it is.")
	addCmd.Flags().StringVarP(&url, "url", "u", "", "Any URL resource associated with this tasks, such as an article.")
	addCmd.Flags().StringVarP(&importFile, "from-file", "f", "", "Import tasks from a file. Can be csv.")
	addCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Go through the steps but do nothing.")
}

// add is the running function for the `add` command.
// It currently has two paths: From file or from CLI.
// From file currently treats the file as it would the CSV database file. This
// means that all the fields must be present, and you can set internal fields
// like UUID. This will change in the future.
// From CLI takes all of the given parameters and autocompletes as best it can
// based on them. It will fill in parent/children relationships, assign UUID,
// and eventually guess priority and size.
// It also operates on failfast, meaning that it will stop immediately if
// anything goes wrong.
func add(cmd *cobra.Command, args []string) error {
	var err error
	var tasks []models.Task

	if err = validate(cmd, args); err != nil {
		return err
	}

	if importFile != "" {
		tasks, err = tasksFromFile()
	} else {
		tasks, err = tasksFromCmd(cmd, args)
	}

	if err != nil {
		return err
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

// validate checks to make sure there aren't any contradictions or out of bound
// fields in the invocation.
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

// tasksFromFile performs the first path in the `add` command. It reads in a CSV
// file that looks exactly like the CSV database file, and creates Task objects.
// TODO: Divine file type, use appropriate storage helpers
func tasksFromFile() ([]models.Task, error) {
	var err error
	var result []models.Task

	f, err := ioutil.ReadFile(importFile)
	if err != nil {
		return []models.Task{}, err
	}
	input := string(f)

	records, err := csv.NewReader(strings.NewReader(input)).ReadAll()
	if err != nil {
		return []models.Task{}, err
	}

	for i, record := range records {
		if i == 0 {
			continue
		}

		newTask, err := storage.TaskFromCsv(record)
		if err != nil {
			return []models.Task{}, err
		}
		result = append(result, newTask)
	}

	return result, nil
}

func tasksFromCmd(cmd *cobra.Command, args []string) ([]models.Task, error) {
	var err error
	var result = make([]models.Task, 1)

	if due_s != "" {
		due, err = time.Parse("2007-01-02", due_s)
		if err != nil {
			return []models.Task{}, err
		}

		due_p, err = ptypes.TimestampProto(due)
		if err != nil {
			return []models.Task{}, err
		}
	}

	result[0] = models.Task{
		Name:     name,
		Size:     uint32(size),
		Tags:     strings.Split(tags, viper.GetString("Delimiter")),
		Added:    ptypes.TimestampNow(),
		Due:      due_p,
		Priority: uint32(priority),
		Url:      url,
	}

	return result, nil
}
