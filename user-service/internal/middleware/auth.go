package middleware

import (
	"net/http"
	"slices"
	"strings"
	"user-service/pkg/jwt"
)

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
			request = PayloadToContext(request, &token.Payload)
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
			request = PayloadToContext(request, &token.Payload)
			nextHandler.ServeHTTP(responseWriter, request)
		},
	)
}
