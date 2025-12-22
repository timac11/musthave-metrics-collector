package middleware

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
)

func (m *Middleware) SetSignatureMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r == nil {
			return
		}

		recorder := &ResponseRecorder{
			ResponseWriter: w,
			body:           &bytes.Buffer{},
			status:         http.StatusOK,
		}

		h.ServeHTTP(recorder, r)

		if m.hashingKey != "" {
			hashInBytes := sha256.Sum256(recorder.body.Bytes())
			signature := hex.EncodeToString(hashInBytes[:])
			r.Header.Set("HashSHA256", signature)
		}

	})
}

func (m *Middleware) CheckSignatureMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r == nil {
			return
		}

		if m.hashingKey != "" {
			bodyBytes, err := io.ReadAll(r.Body)

			if err != nil {
				http.Error(w, "Auth failed", http.StatusInternalServerError)
				return
			}

			r.Body.Close()

			hashInBytes := sha256.Sum256(bodyBytes)
			signature := hex.EncodeToString(hashInBytes[:])

			requestSignature := r.Header.Get("HashSHA256")

			if requestSignature != "" && signature != requestSignature {
				http.Error(w, "Auth failed", http.StatusInternalServerError)
				return
			}

			// reassign body
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		h.ServeHTTP(w, r)
	})
}
