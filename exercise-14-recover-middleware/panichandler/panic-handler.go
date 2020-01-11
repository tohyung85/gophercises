package panichandler

import (
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
)

type panicResponseWriter struct {
	nonPanicCode    int
	wroteHeaderCode bool
	w               http.ResponseWriter
}

func (pw *panicResponseWriter) Header() http.Header {
	return pw.w.Header()
}

func (pw *panicResponseWriter) WriteHeader(code int) {
	if !pw.wroteHeaderCode {
		pw.wroteHeaderCode = true
		pw.nonPanicCode = code
		return
	}
	fmt.Printf("Response code already written as %d vs %d", pw.nonPanicCode, code)
	// debug.PrintStack()
	return
}

func (pw *panicResponseWriter) Write(b []byte) (int, error) {
	if pw.wroteHeaderCode { // Note: Only writes the first written header code
		pw.w.WriteHeader(pw.nonPanicCode)
	}
	return pw.w.Write(b)
}

type PanicHandler struct {
	mux http.Handler
}

func New(mux http.Handler) *PanicHandler {
	return &PanicHandler{mux}
}

func (p *PanicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pw := &panicResponseWriter{w: w, wroteHeaderCode: false}
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Logging error: %s\n", r)
			debug.PrintStack()
			env, _ := os.LookupEnv("ENV")
			if env == "Development" {
				errorString := fmt.Sprintf("%s\n", r) + string(debug.Stack())
				http.Error(w, errorString, http.StatusInternalServerError)
				return
			}
			http.Error(w, "Opps! Something went wrong!", http.StatusInternalServerError)
		}
	}()

	p.mux.ServeHTTP(pw, r)
}
