package middleware

import (
	"bytes"
	"net/http"
)

type ResponseRecorder struct {
	http.ResponseWriter
	body   *bytes.Buffer
	status int
}

func (r *ResponseRecorder) Write(b []byte) (int, error) {
	r.body.Write(b)                  // Record the body
	return r.ResponseWriter.Write(b) // Write to the actual response
}

func (r *ResponseRecorder) WriteHeader(statusCode int) {
	r.status = statusCode                    // Record the status code
	r.ResponseWriter.WriteHeader(statusCode) // Write the actual header
}
