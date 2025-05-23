package jwt

import (
	"encoding/base64"
	"errors"
	"strings"
	"time"
	"user-service/pkg/jwt/jwt_errors"
	"user-service/pkg/jwt/jwt_metadata"
)

type JWTCoder struct {
	alg                    string        // Алгоритм подписи
	secret                 string        // Серкрет
	issuer                 string        // Издательиздатель токена
	audience               []string      // Получатели, которым предназначается данный токен
	expirationTimeDuration time.Duration // Время жизни токена
}

func (c *JWTCoder) NewToken(subject, role string, permissions ...string) *jwt_metadata.Token {
	header := jwt_metadata.Header{
		Alg: c.alg,
		Typ: "JWT",
	}

	payload := jwt_metadata.Payload{
		BasePayload: jwt_metadata.BasePayload{
			Issuer:  c.issuer,
			Subject: subject,
		},
		Role:        role,
		Permissions: permissions,
	}
	payload.SetAudience(c.audience)

	timeNow := time.Now()
	payload.SetExpiration(timeNow, c.expirationTimeDuration)
	payload.SetNotBefore(timeNow, 0)
	payload.SetIssuedAt(timeNow, 0)

	return &jwt_metadata.Token{
		Header:  header,
		Payload: payload,
	}
}

func (c *JWTCoder) Encode(token jwt_metadata.Token) (*string, error) {
	header := &token.Header
	if err := header.Validate(); err != nil {
		return nil, errors.Join(jwt_errors.ErrTokenGeneration, err)
	}
	headerBase64, err := serializeToBase64(header)
	if err != nil {
		return nil, errors.Join(jwt_errors.ErrTokenGeneration, err)
	}

	payload := &token.Payload
	if err := payload.Validate(); err != nil {
		return nil, errors.Join(jwt_errors.ErrTokenGeneration, err)
	}
	payloadBase64, err := serializeToBase64(payload)
	if err != nil {
		return nil, errors.Join(jwt_errors.ErrTokenGeneration, err)
	}

	encodedSignature := generateSignature(headerBase64, payloadBase64, &header.Alg, &c.secret)

	tokenParts := []string{
		*headerBase64,
		*payloadBase64,
		string(base64.RawURLEncoding.EncodeToString(encodedSignature)),
	}

	jwt := strings.Join(tokenParts, ".")

	return &jwt, nil
}

func (c *JWTCoder) Parse(token string) (*jwt_metadata.Token, error) {
	tokenParts := strings.Split(token, ".")
	if len(tokenParts) != 3 {
		return nil, jwt_errors.ErrInvalidTokenFormat
	}

	header, err := parseHeader(&tokenParts[0])
	if err != nil {
		return nil, err
	}

	payload, err := parsePayload(&tokenParts[1])
	if err != nil {
		return nil, err
	}

	if err := validateSignature(&tokenParts[0], &tokenParts[1], &tokenParts[2], &header.Alg, &c.secret); err != nil {
		return nil, err
	}

	return &jwt_metadata.Token{
		Header:  *header,
		Payload: *payload,
	}, nil
}

func NewJWTCoder(alg, secret, issuer string, audience []string, expirationTimeDuration time.Duration) *JWTCoder {
	return &JWTCoder{
		alg:                    alg,
		secret:                 secret,
		issuer:                 issuer,
		audience:               audience,
		expirationTimeDuration: expirationTimeDuration,
	}
}
