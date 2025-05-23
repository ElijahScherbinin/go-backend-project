package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"
	"user-service/pkg/jwt"
	"user-service/pkg/jwt/jwt_helper"
	"user-service/pkg/jwt/jwt_metadata"
)

type ContextKey string

const PayloadContextKey ContextKey = "payload"

const Alg string = "HS256"                                       // TODO: remove
const Secret string = "a-string-secret-at-least-256-bits-long"   // TODO: remove
const Issuer string = "user-service"                             // TODO: remove
var Audience []string = []string{"user-service", "test-service"} // TODO: remove
const ExpirationTimeDuration time.Duration = time.Minute * 15    // TODO: remove

var jwtEncoder *jwt.JWTCoder = jwt.NewJWTCoder(Alg, Secret, Issuer, Audience, ExpirationTimeDuration)

func IsAdminMiddleware(nextHandler http.Handler) http.Handler {
	return http.HandlerFunc(
		func(responseWriter http.ResponseWriter, request *http.Request) {
			token, err := getToken(request)
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
			token, err := getToken(request)
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

func getToken(request *http.Request) (*jwt_metadata.Token, error) {
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

func SavePayload(request *http.Request, payload *jwt_metadata.Payload) *http.Request {
	ctx := context.WithValue(request.Context(), PayloadContextKey, &payload)
	return request.WithContext(ctx)
}

func GetPayload(request *http.Request) (*jwt_metadata.Payload, error) {
	payloadInterface := request.Context().Value(PayloadContextKey)
	if payloadInterface == nil {
		return nil, fmt.Errorf("'%s' not found in context", PayloadContextKey)
	}

	if payload, ok := payloadInterface.(*jwt_metadata.Payload); ok {
		return payload, nil
	}
	return nil, errors.New("invalid payload type in context")
}
