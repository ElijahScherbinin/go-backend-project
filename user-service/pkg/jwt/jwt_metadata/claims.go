package jwt_metadata

import (
	"time"
	"user-service/pkg/jwt/jwt_errors"
)

type Claims struct {
	Issuer         string `json:"iss,omitempty"` // издатель токена
	Subject        string `json:"sub"`           // субъект, которому выдан токен
	Audience       string `json:"aud,omitempty"` // получатели, которым предназначается данный токен
	ExpirationTime int64  `json:"exp"`           // время, когда токен станет невалидным
	NotBefore      int64  `json:"nbf"`           // время, с которого токен должен считаться действительным
	IssuedAt       int64  `json:"iat"`           // время, в которое был выдан токен
	JWTID          string `json:"jti,omitempty"` // уникальный идентификатор токена
}

// SetExpiration устанавливает время истечения токена
func (p *Claims) SetExpiration(duration time.Duration) {
	p.ExpirationTime = time.Now().Add(duration).Unix()
}

// SetNotBefore устанавливает время начала действия токена
func (p *Claims) SetNotBefore(duration time.Duration) {
	p.NotBefore = time.Now().Add(duration).Unix()
}

func (p *Claims) Validate() error {
	if p.ExpirationTime <= 0 {
		return jwt_errors.ErrInvalidClaims
	}
	if p.NotBefore <= 0 {
		return jwt_errors.ErrInvalidClaims
	}
	if p.IssuedAt <= 0 {
		return jwt_errors.ErrInvalidClaims
	}
	if p.ExpirationTime < p.NotBefore {
		return jwt_errors.ErrInvalidTimeRange
	}
	currentTime := time.Now().Unix()
	if currentTime > p.ExpirationTime {
		return jwt_errors.ErrTokenExpired
	}
	if currentTime < p.NotBefore {
		return jwt_errors.ErrNotBeforeError
	}
	return nil
}
