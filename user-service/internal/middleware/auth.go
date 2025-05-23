package middleware

import (
	"net/http"
	"slices"
	"strings"
	"time"
	"user-service/pkg/jwt"
)

const Alg string = "HS256"                                       // TODO: remove
const Secret string = "a-string-secret-at-least-256-bits-long"   // TODO: remove
const Issuer string = "user-service"                             // TODO: remove
var Audience []string = []string{"user-service", "test-service"} // TODO: remove
const ExpirationTimeDuration time.Duration = time.Minute * 15    // TODO: remove

var jwtEncoder *jwt.JWTCoder = jwt.NewJWTCoder(Alg, Secret, Issuer, Audience, ExpirationTimeDuration)

func IsAdminMiddleware(nextHandler http.Handler) http.Handler {
	return http.HandlerFunc(
		func(responseWriter http.ResponseWriter, request *http.Request) {
			token, err := GetToken(request)
			if err != nil {
				http.Error(responseWriter, err.Error(), http.StatusUnauthorized)
				return
			}
			if strings.Compare(token.Payload.Role, "admin") != 0 {
				http.Error(responseWriter, "Доступ запрещен", http.StatusForbidden)
				return
			}
			request = SavePayload(request, &token.Payload)
			nextHandler.ServeHTTP(responseWriter, request)
		},
	)
}

func IsPermitionDeleteMiddleware(nextHandler http.Handler) http.Handler {
	return http.HandlerFunc(
		func(responseWriter http.ResponseWriter, request *http.Request) {
			token, err := GetToken(request)
			if err != nil {
				http.Error(responseWriter, err.Error(), http.StatusUnauthorized)
				return
			}
			if !slices.Contains(token.Payload.Permissions, "delete") {
				http.Error(responseWriter, "Доступ запрещен", http.StatusForbidden)
				return
			}
			request = SavePayload(request, &token.Payload)
			nextHandler.ServeHTTP(responseWriter, request)
		},
	)
}
