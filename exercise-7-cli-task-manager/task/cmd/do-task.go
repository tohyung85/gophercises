package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var doTaskCmd = &cobra.Command{
	Use:   "do",
	Short: "Complete Task",
	Long:  "Set task as complete",
	Args:  cobra.MinimumNArgs(1),
	Run:   doTask,
}

func doTask(cmd *cobra.Command, args []string) {
	taskDone := args[0]
	taskId, err := strconv.Atoi(taskDone)
	if err != nil {
		fmt.Printf("Entry: %s is not an integer!", taskDone)
	}
	TasksStore.FlagComplete(taskId)
	fmt.Printf("You have completed \"%s\" task\n", taskDone)
}
