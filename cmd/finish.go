package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/golang/protobuf/ptypes"
	"github.com/spf13/cobra"
)

var finishFlags struct {
	remove bool
	dryRun bool

	uuid uint64
}

// addCmd represents the add command
var finishCmd = &cobra.Command{
	Use:   "finish",
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
	RunE: finish,
	// TODO:
	// PreRun: validate,
}

func init() {
	TaskerCmd.AddCommand(finishCmd)

	finishCmd.Flags().BoolVarP(&finishFlags.remove, "remove", "r", true, "Also remove the task from the list.")
	finishCmd.Flags().BoolVar(&finishFlags.dryRun, "dry-run", false, "Go through the steps but do nothing.")
}

func finish(cmd *cobra.Command, args []string) error {
	var err error
	defer func() {
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		}
	}()

	if err = validateFinish(cmd, args); err != nil {
		return nil
	}

	task, err := db.GetTask(finishFlags.uuid)
	if err != nil {
		return nil
	}

	oldTask := task

	if task.Finished != nil {
		finished, terr := ptypes.Timestamp(task.Finished)
		if err != nil {
			return terr
		}

		err = fmt.Errorf("Task %d (%s) was finished %s", task.Guid, task.Name, finished.Format("at 15:04 on Mon Jan 2, 2006"))

		if !task.Removed && finishFlags.remove {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
			fmt.Fprintf(os.Stderr, "It has not been removed. Remove it now?\n[Y,n]: ")

			err = nil

			reader := bufio.NewReader(os.Stdin)
			answer, _ := reader.ReadString('\n')

			if answer == "Y" {
				task.Removed = true
				if !finishFlags.dryRun {
					err = db.EditTask(oldTask, task)
				}
			}
		}

		return nil
	}

	task.Finished = ptypes.TimestampNow()
	task.Removed = finishFlags.remove

	if !finishFlags.dryRun {
		err = db.EditTask(oldTask, task)
	}

	return nil
}

// validateFinish checks to make sure there aren't any contradictions or out of
// bound fields in the invocation.
func validateFinish(cmd *cobra.Command, args []string) error {
	// TODO: Implement
	if len(args) == 0 && !cmd.Flag("uuid").Changed {
		return fmt.Errorf("Need to have something to finish!")
	}

	if uuid, err := strconv.ParseUint(args[0], 10, 64); err == nil {
		finishFlags.uuid = uuid
	} else {
		return err
	}

	return nil
}
