package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/tohyung85/gophercises/exercise-13-quiet-hn/handlers/hnhandler"
)

func main() {
	port := 8080
	hnHandler := hnhandler.NewHandler()
	http.Handle("/", hnHandler)
	// fs := http.FileServer(http.Dir(".."))
	// http.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Printf("Starting Server on localhost:%d\n", port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))

}
