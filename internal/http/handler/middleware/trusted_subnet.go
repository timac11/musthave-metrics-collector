package middleware

import (
	"net/http"
	"net/netip"

	"github.com/timac11/musthave-metrics-collector/internal/logger"
)

func (m *Middleware) CheckSubnetMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r == nil {
			return
		}

		if (r.Method == http.MethodPost || r.Method == http.MethodPut) && m.subnet != nil {
			realIP := r.Header.Get("X-Real-IP")
			ip, err := netip.ParseAddr(realIP)

			if err != nil {
				logger.Error("Failed to check client ip", err)
				http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
				return
			}

			if !m.subnet.Contains(ip) {
				logger.Error("Client ip is not in trusted subnet", realIP)
				http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
				return
			}
		}

		h.ServeHTTP(w, r)
	})
}
