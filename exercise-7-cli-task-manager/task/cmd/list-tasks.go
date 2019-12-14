package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var listTasksCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists Tasks",
	Long:  "Lists all outstanding tasks",
	Run:   listTasks,
}

func listTasks(cmd *cobra.Command, args []string) {
	allTasks, err := TasksStore.RetrieveAll()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("You have the following tasks:")
	for k, v := range allTasks {
		if !v.Completed {
			fmt.Printf("%d: %s\n", k, v.Description)
		}
	}
}
