package middleware

import (
	"bytes"
	"io"
	"net/http"
)

func (m *Middleware) DecryptMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r == nil {
			return
		}

		if (r.Method == "POST" || r.Method == "PUT") && m.signingKey != "" {
			bodyBytes, err := io.ReadAll(r.Body)

			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			r.Body.Close()

			if len(bodyBytes) != 0 && m.decoder != nil {
				decodedMessage, err := m.decoder.Decode(bodyBytes)
				if err != nil {
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}

				r.Body = io.NopCloser(bytes.NewBuffer(decodedMessage))
			} else {
				r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}
		}

		h.ServeHTTP(w, r)
	})
}
