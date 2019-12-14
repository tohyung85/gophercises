package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var delTaskCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove Task",
	Long:  "Delete Task from list",
	Args:  cobra.MinimumNArgs(1),
	Run:   delTask,
}

func delTask(cmd *cobra.Command, args []string) {
	taskDone := args[0]
	taskId, err := strconv.Atoi(taskDone)
	if err != nil {
		fmt.Printf("Entry: %s is not an integer!", taskDone)
	}
	TasksStore.DeleteItem(taskId)
	fmt.Printf("You have deleted \"%s\" task\n", taskDone)
}
