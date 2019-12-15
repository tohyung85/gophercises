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
	var taskKey int
	taskDone := args[0]
	taskId, err := strconv.Atoi(taskDone)
	if err != nil {
		fmt.Printf("Entry: %s is not an integer!", taskDone)
	}
	allTasks, err := TasksStore.RetrieveAll()
	if err != nil {
		fmt.Println("Issue getting all tasks")
	}
	taskIterationId := 1
	for k, v := range allTasks {
		if !v.Completed {
			if taskIterationId != taskId {
				taskIterationId++
				continue
			}
			taskKey = k
		}
	}
	TasksStore.DeleteItem(taskKey)
	fmt.Printf("You have deleted \"%s\" task\n", taskDone)
}
