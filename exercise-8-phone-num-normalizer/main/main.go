package main

import (
	// "database/sql" - if not using sqlx

	"flag"
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/jmoiron/sqlx"
	"github.com/tohyung85/gophercises/exercise-8-phone-num-normalizer/normalizer"
	"github.com/tohyung85/gophercises/exercise-8-phone-num-normalizer/store"

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
	dbStylePtr := flag.String("dbstyle", "Standard", "Type of DB package to use: defaults to go standard sql package")
	flag.Parse()
	dbName, found := os.LookupEnv("DB_NAME")
	if !found {
		fmt.Printf("DB Name not found!")
	}
	dbUser, found := os.LookupEnv("DB_USER")
	if !found {
		fmt.Printf("DB User not found!")
	}

	switch *dbStylePtr {
	case "Standard":
		dbStyle = store.SQL
	case "SQLX":
		dbStyle = store.SQLX
	case "GORM":
		dbStyle = store.GORM
	default:
		dbStyle = store.SQL
	}

	dbGen, err := setupDB(dbName, dbUser, *dbStylePtr)
	if err != nil {
		log.Fatal("Error encountered opening database")
	}
	var phoneStore *store.PhoneStore

	switch *dbStylePtr {
	case "Standard":
		dbStyle = store.SQL
		db := dbGen.(*sqlx.DB)
		defer db.Close()
		phoneStore, err = store.NewStore(db)
	case "SQLX":
		dbStyle = store.SQLX
		db := dbGen.(*sqlx.DB)
		defer db.Close()
		phoneStore, err = store.NewStore(db)
	case "GORM":
		db := dbGen.(*gorm.DB)
		defer db.Close()
		phoneStore, err = store.NewGormStore(db)
	default:
		dbStyle = store.SQLX
		db := dbGen.(*sqlx.DB)
		defer db.Close()
		phoneStore, err = store.NewStore(db)
	}

	if err != nil {
		log.Fatal("Error setting up tables")
	}

	listAllEntries(phoneStore)

	cleanUpDb(phoneStore)
}

func setupDB(dbName string, dbUser string, inptStyle string) (interface{}, error) {
	var db interface{}
	var err error
	switch inptStyle {
	case "GORM":
		dbStyle = store.GORM
		connStr := fmt.Sprintf("user=%s dbname=%s sslmode=disable", dbUser, dbName)
		db, err = gorm.Open("postgres", connStr)
		if err != nil {
			log.Fatal("Error encountered opening gorm database")
		}
	default:
		dbStyle = store.SQL
		connStr := fmt.Sprintf("user=%s dbname=%s sslmode=disable", dbUser, dbName)
		db, err = sqlx.Open("postgres", connStr)
		if err != nil {
			log.Fatal("Error encountered opening database")
		}
	}
	return db, err
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
