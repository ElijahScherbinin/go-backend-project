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
	"user-service/pkg/jwt/jwt_errors"
	"user-service/pkg/jwt/jwt_metadata"
)

func serializeToBase64[T jwt_metadata.SerializebleBase64](data *T) (*string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, errors.Join(jwt_errors.ErrConvertToJson, err)
	}
	encodeData := base64.RawURLEncoding.EncodeToString(jsonData)
	return &encodeData, nil
}

func generateSignature(encodedHeader, encodedClaims, algorithm, secret *string) []byte {
	var hash hash.Hash
	bytesSecret := []byte(*secret)

	switch *algorithm {
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

func parseHeader(headerBase64 *string) (*jwt_metadata.Header, error) {
	headerBytes, err := base64.RawURLEncoding.DecodeString(*headerBase64)
	if err != nil {
		return nil, errors.Join(jwt_errors.ErrInvalidHeader, err)
	}
	var header jwt_metadata.Header
	if err := json.Unmarshal(headerBytes, &header); err != nil {
		return nil, errors.Join(jwt_errors.ErrInvalidHeaderFormat, err)
	}
	if err := header.Validate(); err != nil {
		return nil, err
	}
	return &header, nil
}

func parseClaims(claimsBase64 *string) (*jwt_metadata.Claims, error) {
	claimsBytes, err := base64.RawURLEncoding.DecodeString(*claimsBase64)
	if err != nil {
		return nil, errors.Join(jwt_errors.ErrInvalidClaims, err)
	}
	var claims jwt_metadata.Claims
	if err := json.Unmarshal(claimsBytes, &claims); err != nil {
		return nil, errors.Join(jwt_errors.ErrInvalidClaimsFormat, err)
	}
	if err := claims.Validate(); err != nil {
		return nil, err
	}
	return &claims, nil
}

func validateSignature(headerBase64, claimsBase64, signatureBase64, algorithm, secret *string) error {
	verifySignature := generateSignature(headerBase64, claimsBase64, algorithm, secret)
	verifySignatureBase64 := base64.RawURLEncoding.EncodeToString(verifySignature)
	if !hmac.Equal([]byte(*signatureBase64), []byte(verifySignatureBase64)) {
		return jwt_errors.ErrSignatureVerificationFailed
	}
	return nil
}
