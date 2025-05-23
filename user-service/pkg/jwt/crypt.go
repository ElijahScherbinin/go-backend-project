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

func generateSignature(encodedHeader, encodedPayload, algorithm, secret *string) []byte {
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
	buffer.WriteString(*encodedPayload)
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

func parsePayload(payloadBase64 *string) (*jwt_metadata.Payload, error) {
	payloadBytes, err := base64.RawURLEncoding.DecodeString(*payloadBase64)
	if err != nil {
		return nil, errors.Join(jwt_errors.ErrInvalidPayload, err)
	}
	var payload jwt_metadata.Payload
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return nil, errors.Join(jwt_errors.ErrInvalidPayloadFormat, err)
	}
	if err := payload.Validate(); err != nil {
		return nil, err
	}
	return &payload, nil
}

func validateSignature(headerBase64, payloadBase64, signatureBase64, algorithm, secret *string) error {
	verifySignature := generateSignature(headerBase64, payloadBase64, algorithm, secret)
	verifySignatureBase64 := base64.RawURLEncoding.EncodeToString(verifySignature)
	if !hmac.Equal([]byte(*signatureBase64), []byte(verifySignatureBase64)) {
		return jwt_errors.ErrSignatureVerificationFailed
	}
	return nil
}
