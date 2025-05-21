package jwt

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"
)

type contextKey string

const secretKey = "a-string-secret-at-least-256-bits-long" // TODO: remove
const claimsKey contextKey = "claims"

// CustomClaims - пользовательские поля для JWT
type CustomClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	Issuer   string `json:"iss"`
	Expires  int64  `json:"exp"`
}

// NewClaims - создание новых claims
func NewClaims(userID uint, username, role string, expiry time.Duration) *CustomClaims {
	return &CustomClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		Issuer:   "user-service",
		Expires:  time.Now().Add(expiry).Unix(),
	}
}

// generateSignature - создание подписи для JWT
func generateSignature(header, payload string) string {
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(header + "." + payload))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

// GenerateToken - создание JWT токена
func GenerateToken(claims *CustomClaims) (string, error) {
	header := map[string]string{
		"typ": "JWT",
		"alg": "HS256",
	}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", err
	}

	payloadJSON, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	// Кодируем header и payload в base64
	encodedHeader := base64.RawURLEncoding.EncodeToString(headerJSON)
	encodedPayload := base64.RawURLEncoding.EncodeToString(payloadJSON)

	// Генерируем подпись
	encodedSignature := generateSignature(encodedHeader, encodedPayload)

	// Собираем токен
	token := strings.Join([]string{encodedHeader, encodedPayload, encodedSignature}, ".")
	return token, nil
}

// extractToken - извлекает токен из заголовка Authorization
func extractToken(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")
	if len(bearerToken) > 7 && strings.HasPrefix(bearerToken, "Bearer ") {
		return bearerToken[7:]
	}
	return ""
}

// validateToken - валидация JWT токена
func validateToken(tokenString string) (*CustomClaims, error) {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid token format")
	}

	if _, err := base64.RawURLEncoding.DecodeString(parts[0]); err != nil {
		return nil, errors.New("invalid header")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, errors.New("invalid payload")
	}

	var claims CustomClaims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, errors.New("invalid payload format")
	}

	if claims.Expires < time.Now().Unix() {
		return nil, errors.New("token expired")
	}

	// Проверка подписи
	encodedSignature := generateSignature(parts[0], parts[1])
	if !hmac.Equal([]byte(parts[2]), []byte(encodedSignature)) {
		return nil, errors.New("invalid signature")
	}

	return &claims, nil
}

// JWTMiddleware - middleware для проверки JWT токена
func JWTMiddleware(nextHandler http.Handler) http.Handler {
	return http.HandlerFunc(
		func(responseWriter http.ResponseWriter, requerst *http.Request) {
			tokenString := extractToken(requerst)
			if tokenString == "" {
				http.Error(responseWriter, "Missing token", http.StatusUnauthorized)
				return
			}

			claims, err := validateToken(tokenString)
			if err != nil {
				log.Println(err)
				http.Error(responseWriter, err.Error(), http.StatusUnauthorized)
				return
			}

			// Добавляем claims в контекст запроса
			ctx := context.WithValue(requerst.Context(), claimsKey, claims)
			requerst = requerst.WithContext(ctx)

			nextHandler.ServeHTTP(responseWriter, requerst)
		},
	)
}

// GetClaims - извлекает *CustomClaims из *http.Request
func GetClaims(requerst *http.Request) *CustomClaims {
	return requerst.Context().Value("claims").(*CustomClaims)
}
