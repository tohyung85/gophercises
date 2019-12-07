package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/tohyung85/gophercises/exercise-4-link-parser/linkparser"
)

func main() {
	filePathPtr := flag.String("file", "../data/ex1.html", "html file to parse, defaults to data/ex1.html")
	flag.Parse()
	mux := defaultMux()
	port := 8080
	fmt.Printf("Starting server on localhost port: %d\n", port)

	file, err := os.Open(*filePathPtr)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	linkParser := linkparser.NewParser()

	links, err := linkParser.ParseHTML(file)

	for _, link := range links {
		fmt.Print(link)
	}

	http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}
