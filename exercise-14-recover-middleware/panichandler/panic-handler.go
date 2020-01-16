package panichandler

import (
	"bufio"
	"fmt"
	"html/template"
	"net"
	"net/http"
	"regexp"
	"runtime/debug"
)

type panicResponseWriter struct {
	code   int
	writes [][]byte
	w      http.ResponseWriter
}

func (pw *panicResponseWriter) Header() http.Header {
	return pw.w.Header()
}

func (pw *panicResponseWriter) WriteHeader(code int) {
	pw.code = code
	return
}

func (pw *panicResponseWriter) Write(b []byte) (int, error) {
	tmp := make([]byte, len(b))
	copy(tmp, b)
	pw.writes = append(pw.writes, tmp)
	return len(b), nil
}

func (pw *panicResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := pw.w.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("the Response writer does not support Hijacker interface")
	}
	return hijacker.Hijack()
}

func (pw *panicResponseWriter) Flush() {
	flusher, ok := pw.w.(http.Flusher)

	if !ok {
		return
	}
	flusher.Flush()
}

func (pw *panicResponseWriter) flush() {
	if pw.code != 0 {
		pw.w.WriteHeader(pw.code)
	}
	for _, write := range pw.writes {
		_, err := pw.w.Write(write)
		if err != nil {
			fmt.Printf("Error Encountered: %s", err)
			return
		}
	}
}

type PanicHandler struct {
	mux         http.Handler
	environment string
}

func New(mux http.Handler, environment string) *PanicHandler {
	return &PanicHandler{mux, environment}
}

func (p *PanicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	writes := make([][]byte, 0)
	pw := &panicResponseWriter{w: w, writes: writes}
	defer func() {
		if r := recover(); r != nil {
			// fmt.Printf("Logging error: %s\n", r)
			debug.PrintStack()
			if p.environment == "Development" {
				errorString := fmt.Sprintf("%s\n", r) + fmt.Sprintf("<p>%s</p>", getStackString())
				errorHtml := template.HTML(errorString)
				w.WriteHeader(http.StatusInternalServerError)

				t := template.New("debug")
				t, err := t.Parse("<p>{{.}}</p>")
				if err != nil {
					fmt.Println("Error parsing html")
					http.Error(w, errorString, http.StatusInternalServerError)
					return
				}
				t.Execute(w, errorHtml)
				return
			}
			http.Error(w, "Opps! Something went wrong!", http.StatusInternalServerError)
		}
	}()

	p.mux.ServeHTTP(pw, r)
	pw.flush() // Reminder: If panic called in ServeHTTP, this will never be called
}

func getAllLinks() []string {
	stackString := string(debug.Stack())
	pattern := regexp.MustCompile(`(/.*):([0-9]*)`)
	return pattern.FindAllString(stackString, -1)
}

func getStackString() string {
	stackString := string(debug.Stack())
	pattern := regexp.MustCompile(`(/.*):([0-9]*)`)
	newString := pattern.ReplaceAllString(stackString, `<br/><a href='/debug?file=$1&line=$2'>$1:$2</a>`)
	return newString
}

func parseStackTrace(stackTrace string) string {
	return stackTrace
}
