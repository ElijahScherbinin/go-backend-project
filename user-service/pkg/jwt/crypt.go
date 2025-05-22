package jwt

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"errors"
	"hash"
	"strings"
	"user-service/pkg/jwt/jwt_errors"
	"user-service/pkg/jwt/jwt_metadata"
)

type JWTEncoder struct {
	header jwt_metadata.Header
	secret string
}

func serializeToBase64[T jwt_metadata.SerializebleBase64](data T) (*string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, errors.Join(jwt_errors.ErrConvertToJson, err)
	}
	encodeData := base64.RawURLEncoding.EncodeToString(jsonData)
	return &encodeData, nil
}

func generateSignature(encodedHeader, encodedClaims *string, secret, algorithm string) []byte {
	var hash hash.Hash
	bytesSecret := []byte(secret)

	switch algorithm {
	case "HS256":
		hash = hmac.New(sha256.New, bytesSecret)
	case "HS384":
		hash = hmac.New(sha512.New384, bytesSecret)
	case "HS512":
		hash = hmac.New(sha512.New, bytesSecret)
	default:
		panic(jwt_errors.ErrUnsupportedAlgorithm)
	}

	var buffer bytes.Buffer
	buffer.WriteString(*encodedHeader)
	buffer.WriteRune('.')
	buffer.WriteString(*encodedClaims)
	buffer.WriteTo(hash)

	return hash.Sum(nil)
}

func (e JWTEncoder) Encode(claims jwt_metadata.Claims) (*string, error) {
	if err := e.header.Validate(); err != nil {
		return nil, errors.Join(jwt_errors.ErrTokenGeneration, err)
	}
	headerBase64, err := serializeToBase64(e.header)
	if err != nil {
		return nil, errors.Join(jwt_errors.ErrTokenGeneration, err)
	}

	if err := claims.Validate(); err != nil {
		return nil, errors.Join(jwt_errors.ErrTokenGeneration, err)
	}
	claimsBase64, err := serializeToBase64(claims)
	if err != nil {
		return nil, errors.Join(jwt_errors.ErrTokenGeneration, err)
	}

	encodedSignature := generateSignature(headerBase64, claimsBase64, e.secret, e.header.Alg)

	tokenParts := []string{
		*headerBase64,
		*claimsBase64,
		string(base64.RawURLEncoding.EncodeToString(encodedSignature)),
	}

	token := strings.Join(tokenParts, ".")

	return &token, nil
}

func (e JWTEncoder) ExtractClaims(token string) (*jwt_metadata.Claims, error) {
	tokenParts := strings.Split(token, ".")
	if len(tokenParts) != 3 {
		return nil, jwt_errors.ErrInvalidTokenFormat
	}

	headerBase64 := tokenParts[0]
	headerByte, err := base64.RawURLEncoding.DecodeString(headerBase64)
	if err != nil {
		return nil, errors.Join(jwt_errors.ErrInvalidHeader, err)
	}
	var header jwt_metadata.Header
	if err := json.Unmarshal(headerByte, &header); err != nil {
		return nil, errors.Join(jwt_errors.ErrInvalidHeaderFormat, err)
	}
	if err := header.Validate(); err != nil {
		return nil, err
	}

	claimsBase64 := tokenParts[1]
	claimsByte, err := base64.RawURLEncoding.DecodeString(claimsBase64)
	if err != nil {
		return nil, errors.Join(jwt_errors.ErrInvalidClaims, err)
	}
	var claims jwt_metadata.Claims
	if err := json.Unmarshal(claimsByte, &claims); err != nil {
		return nil, errors.Join(jwt_errors.ErrInvalidClaimsFormat, err)
	}
	if err := claims.Validate(); err != nil {
		return nil, err
	}

	signatureBase64 := tokenParts[2]
	newSignatureBase64 := base64.RawURLEncoding.EncodeToString(
		generateSignature(&headerBase64, &claimsBase64, e.secret, e.header.Alg),
	)
	if !hmac.Equal([]byte(signatureBase64), []byte(newSignatureBase64)) {
		return nil, jwt_errors.ErrSignatureVerificationFailed
	}

	return &claims, nil
}

func NewJWTEncoder(alg, secret string) *JWTEncoder {
	return &JWTEncoder{
		header: *jwt_metadata.NewHeader(alg),
		secret: secret,
	}
}
