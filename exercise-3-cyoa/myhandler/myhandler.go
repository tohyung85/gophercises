package myhandler

import (
	"fmt"
	"net/http"
)

type StoryHandler struct{}

func (*StoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Print("Handling request", r.RequestURI)
	fmt.Fprintln(w, "Hello World!")
	return
}
