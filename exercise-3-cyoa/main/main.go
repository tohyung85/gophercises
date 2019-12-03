package main

import (
	"fmt"
	"github.com/tohyung85/gophercises/exercise-3-cyoa/myhandler"
	"net/http"
)

func main() {
	myHandler := new(myhandler.StoryHandler)

	fmt.Println("Starting Server on localhost:8080")

	http.ListenAndServe(":8080", myHandler)
}
