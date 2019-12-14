package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tohyung85/gophercises/exercise-7-cli-task-manager/task/store"
)

var rootCmd = &cobra.Command{
	Use:   "task",
	Short: "Terminal based Cli Application",
	Long:  "task is a CLI for managing your TODOs.",
}

var TasksStore *store.BoltStore

func Execute() error {
	TasksStore = store.NewStore()
	defer TasksStore.DB.Close()

	rootCmd.AddCommand(addTaskCmd)
	rootCmd.AddCommand(listTasksCmd)
	rootCmd.AddCommand(doTaskCmd)
	rootCmd.AddCommand(delTaskCmd)
	rootCmd.AddCommand(completedTasksCmd)
	return rootCmd.Execute()
}
