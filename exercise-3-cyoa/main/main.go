package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/tohyung85/gophercises/exercise-3-cyoa/myhandler"
)

func main() {
	port := 8080
	platform := flag.String("platform", "web", "Selection to run on 'console' or 'web'. Defaults to 'web'")
	flag.Parse()

	if *platform == "console" {
		fmt.Println("Runnning on Console")
		myHandler := myhandler.NewConsoleHandler()
		myHandler.Start()
	} else {
		myHandler := myhandler.NewWebHandler()
		http.Handle("/", myHandler)
		fs := http.FileServer(http.Dir(".."))
		http.Handle("/static/", http.StripPrefix("/static/", fs))

		fmt.Printf("Starting Server on localhost:%d", port)

		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
	}
}
