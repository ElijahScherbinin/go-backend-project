package jwt_helper

import (
	"net/http"
	"strings"
	"user-service/pkg/jwt/jwt_errors"
)

func ExtractToken(request *http.Request) (string, error) {
	bearerToken := request.Header.Get("Authorization")
	if len(bearerToken) > 7 && strings.HasPrefix(bearerToken, "Bearer ") {
		return bearerToken[7:], nil
	}
	return "", jwt_errors.ErrExtractTokenIsEmpty
}
