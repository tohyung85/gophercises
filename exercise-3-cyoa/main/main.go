package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/tohyung85/gophercises/exercise-3-cyoa/myhandler"
)

func main() {
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

		fmt.Println("Starting Server on localhost:8080")

		http.ListenAndServe(":8080", nil)
	}
}
