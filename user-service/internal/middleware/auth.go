package middleware

import (
	"net/http"
	"slices"
	"time"
	"user-service/pkg/jwt"
	"user-service/pkg/jwt/jwt_helper"
	"user-service/pkg/jwt/jwt_metadata"
)

const Alg string = "HS256"                                       // TODO: remove
const Secret string = "a-string-secret-at-least-256-bits-long"   // TODO: remove
const Issuer string = "user-service"                             // TODO: remove
var Audience []string = []string{"user-service", "test-service"} // TODO: remove
const ExpirationTimeDuration time.Duration = time.Minute * 15    // TODO: remove

var jwtEncoder *jwt.JWTCoder = jwt.NewJWTCoder(Alg, Secret, Issuer, Audience, ExpirationTimeDuration)

// TODO: Доработать
func IsAdminMiddleware(nextHandler http.Handler) http.Handler {
	return http.HandlerFunc(
		func(responseWriter http.ResponseWriter, requerst *http.Request) {
			token, err := GetToken(requerst)
			if err != nil {
				http.Error(responseWriter, err.Error(), http.StatusUnauthorized)
				return
			}
			if !slices.Contains(token.Claims.Permissions, "delete") {
				http.Error(responseWriter, "Доступ запрещен", http.StatusForbidden)
				return
			}
			nextHandler.ServeHTTP(responseWriter, requerst)
		},
	)
}

func GetToken(request *http.Request) (*jwt_metadata.Token, error) {
	tokenString, err := jwt_helper.ExtractToken(request)
	if err != nil {
		return nil, err
	}
	token, err := jwtEncoder.Parse(tokenString)
	if err != nil {
		return nil, err
	}
	return token, nil
}
