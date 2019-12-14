package main

import (
	"fmt"
	"log"

	"github.com/boltdb/bolt"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/tohyung85/gophercises/exercise-7-cli-task-manager/task/cmd"
)

func main() {
	dir, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}
	dbFile := fmt.Sprintf("%s/Go/src/github.com/tohyung85/gophercises/exercise-7-cli-task-manager/task/store/tasks.db", dir)
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	cmd.Execute(db)
}
