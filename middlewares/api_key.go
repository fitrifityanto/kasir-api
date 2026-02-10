package middlewares

import (
	"kasir-api/utils"
	"net/http"
)

// func (api_key) func handler http.handler
func APIKey(validApiKey string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// before function running
			apiKey := r.Header.Get("X-API-Key")

			if apiKey == "" {
				utils.MessageResponse(w, http.StatusUnauthorized, "API Key required", nil)

				return
			}

			if apiKey != validApiKey {
				utils.MessageResponse(w, http.StatusUnauthorized, "Invalid API Key", nil)

				return
			}

			next(w, r)

			// after funciton running
		}
	}
}
