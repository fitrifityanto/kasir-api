package middlewares

import (
	"net/http"
	"strings"

	"github.com/spf13/viper"
)

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		allowedOrigins := viper.GetString("FRONTEND_URL")

		if isAllowed(origin, allowedOrigins) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}

		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-API-Key")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func isAllowed(origin string, allowedString string) bool {
	if origin == "" {
		return false
	}

	origins := strings.SplitSeq(allowedString, ",")
	for o := range origins {
		if strings.TrimSpace(o) == origin {
			return true
		}
	}
	return false
}
