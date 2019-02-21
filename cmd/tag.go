package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var tagFlags struct {
	list bool
}

// tagCmd represents the get command
var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: "A longer description that spans ",
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
	RunE: tag,
}

func init() {
	TaskerCmd.AddCommand(tagCmd)

	tagCmd.Flags().BoolVarP(&tagFlags.list, "list", "l", false, "List all tags in use.")
}

func tag(cmd *cobra.Command, args []string) error {
	// If no options specified, and no target/tag pair, just list
	if tagFlags.list == true || len(args) == 0 {

		tags := db.GetAllTags()
		for _, tag := range tags {
			fmt.Printf("%s, ", tag)
		}

		return nil
	}

	return nil
}
