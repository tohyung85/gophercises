package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/tohyung85/gophercises/exercise-14-recover-middleware/panichandler"
	"log"
	"net/http"
	"os"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Print("No env file found")
	}
}

func main() {
	port := 8080
	environment, found := os.LookupEnv("ENV")
	if !found {
		environment = "Production"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", helloFunc)
	mux.HandleFunc("/panic-demo", panicFunc)
	mux.HandleFunc("/panic-after", panicAfterFunc)

	ph := panichandler.New(mux, environment)

	fmt.Printf("Starting %s server on localhost %d", environment, port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), ph))
}

func helloFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello There!")
}

func panicFunc(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusAccepted)
	panic("oh dear")
}

func panicAfterFunc(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintln(w, "Hello There")
	panic("Oh dear, panicked after")
}
