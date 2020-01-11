package panichandler

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
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
	pw.writes = append(pw.writes, b)
	fmt.Printf("Writing: %s\n", string(b))
	return len(b), nil
}

func (pw *panicResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := pw.w.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("the Response writer does not suppoer Hijacker interface")
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
	pw := &panicResponseWriter{w: w}
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Logging error: %s\n", r)
			debug.PrintStack()
			if p.environment == "Development" {
				errorString := fmt.Sprintf("%s\n", r) + string(debug.Stack())
				http.Error(w, errorString, http.StatusInternalServerError)
				return
			}
			http.Error(w, "Opps! Something went wrong!", http.StatusInternalServerError)
		}
	}()

	p.mux.ServeHTTP(pw, r)
	pw.flush() // Reminder: If panic called in ServeHTTP, this will never be called
}
