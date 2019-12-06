package myhandler

import (
	"bufio"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Page struct {
	Title   string    `json:"title"`
	Story   [1]string `json:"story"`
	Options []option  `json:"options"`
}

type option struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

type StoryHandler struct {
	templ  *template.Template
	UrlMap map[string]Page
}

type ConsoleHandler struct {
	UrlMap map[string]Page
}

func NewConsoleHandler() *ConsoleHandler {
	urlMap, err := setupStories()
	if err != nil {
		panic(err)
	}

	return &ConsoleHandler{
		UrlMap: urlMap,
	}
}

func (c *ConsoleHandler) Start() {
	currPage := c.UrlMap["intro"]
	inpOptMap := mapInputToArc(currPage.Options)

	printPage(&currPage)

	for len(currPage.Options) > 0 {
		reader := bufio.NewReader(os.Stdin)
		inp, _ := reader.ReadString('\n')
		inp = strings.TrimSpace(inp)
		if inpNo, err := strconv.Atoi(inp); err == nil {
			arc, inMap := inpOptMap[inpNo]
			if inMap {
				currPage = c.UrlMap[arc]
				inpOptMap = mapInputToArc(currPage.Options)
				printPage(&currPage)
			}
		}
	}
}

func mapInputToArc(opts []option) map[int]string {
	inpOptMap := make(map[int]string)
	for idx, opt := range opts {
		inpOptMap[idx+1] = opt.Arc
	}

	return inpOptMap
}

func printPage(page *Page) {
	fmt.Println(page.Title)
	fmt.Println("--------------------------------------------------------")
	fmt.Println(page.Story[0])
	fmt.Println("--------------------------------------------------------")

	for idx, opt := range page.Options {
		fmt.Printf("%d: %s\n", idx+1, opt.Text)
	}
}

func NewWebHandler() *StoryHandler {
	filePath, _ := filepath.Abs("../story-template.html")
	tmpl, err := template.ParseFiles(filePath)
	if err != nil {
		panic(err)
	}

	urlMap, err := setupStories()
	if err != nil {
		panic(err)
	}

	return &StoryHandler{
		templ:  tmpl,
		UrlMap: urlMap,
	}
}

func (s *StoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handling request", r.RequestURI)
	path := strings.Replace(r.RequestURI, "/", "", -1)
	data, inMap := s.UrlMap[path]
	if !inMap {
		data = s.UrlMap["intro"]
	}
	s.templ.Execute(w, data)

	return
}

func setupStories() (map[string]Page, error) {
	jsonPath, _ := filepath.Abs("../gopher.json")
	jsonStr, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		panic(err)
	}
	urlMap := make(map[string]Page)

	err = json.Unmarshal(jsonStr, &urlMap)
	if err != nil {
		panic(err)
	}

	return urlMap, err
}
