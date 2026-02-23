package middleware

import (
	"bytes"
	"io"
	"net/http"

	"github.com/timac11/musthave-metrics-collector/internal/common/util"
	"github.com/timac11/musthave-metrics-collector/internal/logger"
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

		if m.signingKey != "" {
			signature, err := util.CalculateSignature(recorder.body.Bytes(), m.signingKey)

			if err != nil {
				logger.Error("Failed to calculate signature", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			r.Header.Set("HashSHA256", signature)
		}

	})
}

func (m *Middleware) CheckSignatureMiddleware(h http.Handler) http.Handler {
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

			if len(bodyBytes) != 0 {
				signature, err := util.CalculateSignature(bodyBytes, m.signingKey)

				if err != nil {
					logger.Error("Failed to calculate signature", err)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}

				requestSignature := r.Header.Get("HashSHA256")

				if requestSignature != "" && signature != requestSignature {
					http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
					return
				}
			}

			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		h.ServeHTTP(w, r)
	})
}
