package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var completedTasksCmd = &cobra.Command{
	Use:   "completed",
	Short: "Lists Completed Tasks",
	Long:  "Lists all completed tasks",
	Run:   listCompletedTasks,
}

func listCompletedTasks(cmd *cobra.Command, args []string) {
	allTasks, err := TasksStore.RetrieveAll()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("You have completed the following tasks:")
	for k, v := range allTasks {
		if v.Completed {
			fmt.Printf("%d: %s\n", k, v.Description)
		}
	}
}
