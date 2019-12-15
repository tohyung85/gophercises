package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"strings"
)

var addTaskCmd = &cobra.Command{
	Use:   "add",
	Short: "Add Task",
	Long:  "Adds a task to our list",
	Args:  cobra.MinimumNArgs(1),
	Run:   addTask,
}

func addTask(cmd *cobra.Command, args []string) {
	taskToAdd := strings.Join(args, " ")
	TasksStore.AddToStore(taskToAdd)
	fmt.Printf("Adding a Task: %s\n", taskToAdd)
}
