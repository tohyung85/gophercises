package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/tohyung85/gophercises/urlshort"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Print("No .env file found")
	}
}

//main
func main() {
	mux := defaultMux()
	ymlFilePath := flag.String("ymlPath", "./path-maps.yaml", "path to yaml file")
	jsonFilePath := flag.String("jsonPath", "./path-maps.json", "path the json file")
	flag.Parse()

	yaml, err := ioutil.ReadFile(*ymlFilePath)
	if err != nil {
		log.Fatal(err)
	}

	jsn, err := ioutil.ReadFile(*jsonFilePath)
	if err != nil {
		log.Fatal(err)
	}

	dbName, exists := os.LookupEnv("DB_NAME")
	if !exists {
		fmt.Print("DB Name not found")
	}

	dbUser, exists := os.LookupEnv("DB_USER")
	if !exists {
		fmt.Print("DB User not found")
	}

	connStr := fmt.Sprintf("user=%s dbname=%s sslmode=disable", dbUser, dbName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	yamlHandler, err := urlshort.YAMLHandler(yaml, mapHandler)
	if err != nil {
		panic(err)
	}

	jsonHandler, err := urlshort.JSONHandler(jsn, yamlHandler)
	if err != nil {
		panic(err)
	}

	dbHandler, err := urlshort.DBHandler(db, jsonHandler)

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", dbHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
