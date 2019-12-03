package myhandler

import (
	"fmt"
	"html/template"
	"net/http"
)

type Page struct {
	Title   string
	Story   string
	Options []option
}

type option struct {
	Text string
	Arc  string
}

type StoryHandler struct {
	templ *template.Template
}

func NewHandler() *StoryHandler {
	tmpl, err := template.ParseFiles("../story-template.html")
	if err != nil {
		panic(err)
	}

	return &StoryHandler{
		templ: tmpl,
	}
}

func (s *StoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Print("Handling request", r.RequestURI)
	opts := make([]option, 2)
	opts[0] = option{Text: "go left", Arc: "going left"}
	opts[1] = option{Text: "go right", Arc: "going right"}
	data := Page{
		Title:   "New Page",
		Story:   "yak yakity yak",
		Options: opts,
	}
	s.templ.Execute(w, data)

	return
}

func readJSON() {

}
