package internal

import (
	"log/slog"

	authclient "github.com/Bitummit/go_auth/pkg/auth_client"
	grpcConfig "github.com/Bitummit/go_auth/pkg/config"
	"github.com/Bitummit/blog_api_golang/pkg/logger"
	"github.com/Bitummit/blog_api_golang/pkg/utils"

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
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request ) {
			token := r.Header.Get("Authorization")
			if token == "" {
				log.Error("Token is empty")
				w.WriteHeader(http.StatusUnauthorized)
				render.JSON(w, r, utils.Error("unauthorized"))
				return
			}
			
			client, err := authclient.NewClient(log, grpcConfig.InitConfig())
			if err != nil {
				log.Error("Error starting grpc auth client", logger.Err(err))
				return
			}

			response, err := client.CheckToken(token)
			defer client.Conn.Close()
			if err != nil || response.Status != "OK" {
				log.Error("invalid token")
				w.WriteHeader(http.StatusInternalServerError)
				render.JSON(w, r, utils.Error("invalid token"))
				return
			} else {
				log.Info("Valid token")
				next.ServeHTTP(w, r)
			}
		})
	}
}