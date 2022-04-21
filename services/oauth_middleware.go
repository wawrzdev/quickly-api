package services

import (
	"fmt"
	"github.com/Nerzal/gocloak/v11"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"strings"
)

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	Title       string `json:"Title"`
	Description string `json:"Description"`
}

type OAuthClientMiddleware struct {
	client       gocloak.GoCloak
	clientId     string
	clientSecret string
	realm        string
}

func NewOauthClientMiddleware() *OAuthClientMiddleware {
	clientId := os.Getenv("AUTH_CLIENT_ID")
	clientSecret := os.Getenv("AUTH_CLIENT_SECRET")
	realm := os.Getenv("AUTH_REALM")
	hostname := os.Getenv("AUTH_HOST")
	return &OAuthClientMiddleware{client: gocloak.NewClient(hostname), clientId: clientId, clientSecret: clientSecret, realm: realm}
}

func (m *OAuthClientMiddleware) Func() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			authHeader := r.Header.Get("Authorization")

			if len(authHeader) < 1 {
				http.Error(w, "Unauthorized Request", http.StatusUnauthorized)
				return
			}

			accessToken := strings.Split(authHeader, " ")[1]

			rptResult, err := m.client.RetrospectToken(r.Context(), accessToken, m.clientId, m.clientSecret, m.realm)

			if err != nil {
				message := fmt.Sprintf("Bad Request: %s", err)
				http.Error(w, message, http.StatusBadRequest)
				return
			}

			isTokenValid := *rptResult.Active

			if !isTokenValid {
				http.Error(w, "Unauthorized Request", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
