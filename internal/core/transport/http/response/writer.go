package core_http_response

import "net/http"

var (
	StatusCodeUnintialized = -1
)

type ResponseWriter struct {
	http.ResponseWriter

	statusCode int
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
		statusCode:     StatusCodeUnintialized,
	}
}

func (rw *ResponseWriter) WriteHeader(statusCode int) {
	rw.ResponseWriter.WriteHeader(statusCode)
	rw.statusCode = statusCode
}

func (rw *ResponseWriter) GetStatusCodeOrPanic() int {
	if rw.statusCode == StatusCodeUnintialized {
		panic("status code not set")
	}
	return rw.statusCode
}
