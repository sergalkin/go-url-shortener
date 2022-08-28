package middleware

import (
	"net"
	"net/http"

	"github.com/sergalkin/go-url-shortener.git/internal/app/config"
)

func TrustedSubnet(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if config.TrustedSubnet() == "" {
			http.Error(writer, "Trusted subnet is not defined. Forbidden.", http.StatusForbidden)
			return
		}

		xRealIP := request.Header.Get("X-Real-IP")
		if xRealIP == "" {
			http.Error(writer, "X-Real-IP is not present.", http.StatusBadRequest)
			return
		}

		ip := net.ParseIP(xRealIP)
		_, subnet, err := net.ParseCIDR(config.TrustedSubnet())
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		if !subnet.Contains(ip) {
			http.Error(writer, "This IP address does not belongs to defined trusted subnet.", http.StatusForbidden)
			return
		}

		next.ServeHTTP(writer, request)
	})
}
