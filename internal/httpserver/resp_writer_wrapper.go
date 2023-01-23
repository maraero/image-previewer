package httpserver

import (
	"net/http"
)

type responseWriterWrapper struct {
	w          *http.ResponseWriter
	statusCode *int
}

func newResponseWriterWrapper(w http.ResponseWriter) responseWriterWrapper {
	statusCode := 200
	return responseWriterWrapper{w: &w, statusCode: &statusCode}
}

// Header overwrites the http.ResponseWriter Header() func.
func (rww responseWriterWrapper) Header() http.Header {
	return (*rww.w).Header()
}

func (rww responseWriterWrapper) Write(buf []byte) (int, error) {
	return (*rww.w).Write(buf)
}

// WriteHeader overwrites the http.ResponseWriter WriteHeader() func.
func (rww responseWriterWrapper) WriteHeader(statusCode int) {
	(*rww.statusCode) = statusCode
	(*rww.w).WriteHeader(statusCode)
}
