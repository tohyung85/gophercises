package cmd

import (
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/spf13/cobra"
	"github.com/tohyung85/gophercises/exercise-7-cli-task-manager/task/store"
)

var rootCmd = &cobra.Command{
	Use:   "task",
	Short: "Terminal based Cli Application",
	Long:  "task is a CLI for managing your TODOs.",
}

var TasksStore *store.BoltStore

func Execute(db *bolt.DB) error {
	setupStore(db)
	rootCmd.AddCommand(addTaskCmd)
	rootCmd.AddCommand(listTasksCmd)
	rootCmd.AddCommand(doTaskCmd)
	rootCmd.AddCommand(delTaskCmd)
	rootCmd.AddCommand(completedTasksCmd)
	return rootCmd.Execute()
}

func setupStore(db *bolt.DB) error {
	TasksStore = store.NewStore(db)
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("DB"))
		if err != nil {
			return fmt.Errorf("Could not create bucket: %v", err)
		}
		return nil
	})
	return err
}
