package utils

import (
	"blog_api/pkg/logger"
	"bytes"
	"encoding/json"
	"log/slog"

	"net/http"

	"github.com/go-chi/render"
)

func SetJSONContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

type CheckTokenResponse struct{
	Status string `json:"status"`
	Error string `json:"error,omitempty"`
}


func CheckTokenMiddleware(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
// func CheckTokenMiddleware(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request ) {
			token := r.Header.Get("Authorization")
			if token == "" {
				log.Error("Token is empty")
				w.WriteHeader(http.StatusInternalServerError)
				render.JSON(w, r, Error("unauthorized"))
				return
			}
			postBody, _ := json.Marshal(map[string]string{
				"token":  token,
			})
			bodyBytes := bytes.NewBuffer(postBody)
			resp, err := http.Post("http://localhost:8001/token/", "application/json", bodyBytes)
			if err != nil {
				log.Error("Error checking token", logger.Err(err))
				w.WriteHeader(http.StatusInternalServerError)
				render.JSON(w, r, Error("internal server error"))
				return
			}

			var response CheckTokenResponse
			_ = render.DecodeJSON(resp.Body, &response)
			if response.Status == "Error" {
				log.Error("invalid token")
				w.WriteHeader(http.StatusUnauthorized)
				render.JSON(w, r, Error("invalid token"))
				return
			} else {
			next.ServeHTTP(w, r)
			}
		})
	}
}