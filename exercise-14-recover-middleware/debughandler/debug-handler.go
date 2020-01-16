package debughandler

import (
	"fmt"
	_ "github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"io/ioutil"
	"net/http"
	"strconv"
)

type DebugHandler struct {
}

func (dh *DebugHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fileName := r.FormValue("file")
	errorLine := r.FormValue("line")
	fmt.Println("Error in line", errorLine)
	file := "/" + fileName
	contents, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	code := string(contents)

	line, err := strconv.Atoi(errorLine)

	err = formatCode(&w, code, line)

	if err != nil {
		fmt.Println(err)
		fmt.Fprintf(w, code)
	}
}

func formatCode(w *http.ResponseWriter, code string, highlightLine int) error {
	lexer := lexers.Get("go")
	if lexer == nil {
		lexer = lexers.Fallback
	}
	highlightRange := [][2]int{[2]int{highlightLine, highlightLine}}
	formatter := html.New(html.HighlightLines(highlightRange), html.Standalone(true), html.WithLineNumbers(true))
	style := styles.Get("dracula")
	if style == nil {
		style = styles.Fallback
	}
	iterator, err := lexer.Tokenise(nil, code)
	if err != nil {
		return err
	}

	err = formatter.Format(*w, style, iterator)
	if err != nil {
		return err
	}
	fmt.Println("Done")
	return nil
}

func New() *DebugHandler {
	return &DebugHandler{}
}
