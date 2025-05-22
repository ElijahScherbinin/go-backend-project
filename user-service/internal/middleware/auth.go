package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"user-service/pkg/jwt"
	"user-service/pkg/jwt/jwt_helper"
	"user-service/pkg/jwt/jwt_metadata"
)

type contextKey string

var ClaimsKey contextKey = "claims"

const Alg string = "HS256"                                     // TODO: remove
const Secret string = "a-string-secret-at-least-256-bits-long" // TODO: remove

var jwtEncoder *jwt.JWTEncoder = jwt.NewJWTEncoder(Alg, Secret)

// JWTMiddleware - middleware для проверки JWT токена
func JWTMiddleware(nextHandler http.Handler) http.Handler {
	return http.HandlerFunc(
		func(responseWriter http.ResponseWriter, requerst *http.Request) {
			tokenString, err := jwt_helper.ExtractToken(requerst)
			if err != nil {
				log.Println("JWTMiddleware:", err)
				http.Error(responseWriter, err.Error(), http.StatusUnauthorized)
				return
			}

			claims, err := jwtEncoder.ExtractClaims(tokenString)
			if err != nil {
				log.Println("JWTMiddleware:", err)
				http.Error(responseWriter, err.Error(), http.StatusUnauthorized)
				return
			}

			// Добавляем claims в контекст запроса
			ctx := context.WithValue(requerst.Context(), ClaimsKey, claims)
			requerst = requerst.WithContext(ctx)

			nextHandler.ServeHTTP(responseWriter, requerst)
		},
	)
}

// TODO: Доработать
func IsAdminMiddleware(nextHandler http.Handler) http.Handler {
	return JWTMiddleware(http.HandlerFunc(
		func(responseWriter http.ResponseWriter, requerst *http.Request) {
			claims, err := GetClaims(requerst)
			if err != nil {
				http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
			}
			if claims.Subject != "admin" {
				http.Error(responseWriter, "Доступ запрещен", http.StatusForbidden)
			}
			nextHandler.ServeHTTP(responseWriter, requerst)
		},
	))
}

func GetClaims(request *http.Request) (*jwt_metadata.Claims, error) {
	claimsInterface := request.Context().Value(ClaimsKey)
	if claimsInterface == nil {
		return nil, fmt.Errorf("claims not found in context")
	}

	if claims, ok := claimsInterface.(*jwt_metadata.Claims); ok {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid claims type in context")
}
