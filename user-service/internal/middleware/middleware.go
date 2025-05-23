package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"user-service/pkg/jwt/jwt_helper"
	"user-service/pkg/jwt/jwt_metadata"
)

type ContextKey string

const PayloadContextKey ContextKey = "payload"

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
