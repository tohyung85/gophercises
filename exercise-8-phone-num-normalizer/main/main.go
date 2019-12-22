package main

import (
	"database/sql"
	"fmt"
	"github.com/tohyung85/gophercises/exercise-8-phone-num-normalizer/normalizer"
	"github.com/tohyung85/gophercises/exercise-8-phone-num-normalizer/store"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var dbStyle store.DbStyle = store.SQL

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Print("No env file found")
	}
}
func main() {
	dbName, found := os.LookupEnv("DB_NAME")
	if !found {
		fmt.Printf("DB Name not found!")
	}
	dbUser, found := os.LookupEnv("DB_USER")
	if !found {
		fmt.Printf("DB User not found!")
	}

	connStr := fmt.Sprintf("user=%s dbname=%s sslmode=disable", dbUser, dbName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error encountered opening database")
	}
	defer db.Close()

	phoneStore, err := store.NewStore(db)
	if err != nil {
		log.Fatal("Error setting up tables")
	}

	listAllEntries(phoneStore)

	cleanUpDb(phoneStore)

}

func cleanUpDb(ps *store.PhoneStore) {
	allEntries, err := ps.GetEntries(dbStyle)
	if err != nil {
		fmt.Printf("Error getting entries: %s\n", err)
		return
	}
	phoneNoMap := make(map[string]bool)
	for _, ent := range allEntries {
		normalized, err := normalizer.NormalizeNumber(ent.Number)
		if err != nil {
			fmt.Printf("Error: Unable to normalize value: %s\n", ent.Number)
			continue
		}
		_, inMap := phoneNoMap[normalized]
		if inMap {
			ps.DeleteEntry(ent.Id, dbStyle)
			continue
		}
		phoneNoMap[normalized] = true
		ps.UpdateEntry(ent.Id, normalized, dbStyle)
	}
}

func setUpMap(pArr []string) map[string]bool {
	phoneNoMap := make(map[string]bool)
	for _, p := range pArr {
		phoneNoMap[p] = true
	}
	return phoneNoMap
}

func listAllEntries(ps *store.PhoneStore) {
	allEntries, err := ps.GetEntries(dbStyle)
	if err != nil {
		fmt.Printf("Error getting entries: %s\n", err)
		return
	}
	fmt.Printf("Current Entries:\n")
	for _, ent := range allEntries {
		fmt.Printf("%s\n", ent.Number)
	}
}
