package jwt_metadata

import (
	"slices"
	"user-service/pkg/jwt/jwt_errors"
)

// SupportedAlgorithms содержит список поддерживаемых алгоритмов
var SupportedAlgorithms = []string{
	"HS256",
	"HS384",
	"HS512",
}

type Header struct {
	Alg string `json:"alg"` // алгоритм подписи
	Typ string `json:"typ"` // тип токена
}

func (h *Header) Validate() error {
	if h.Typ != "JWT" {
		return jwt_errors.ErrUnsupportedTypeToken
	}
	if h.Alg == "" {
		return jwt_errors.ErrMissingAlgorithm
	}
	if !slices.Contains(SupportedAlgorithms, h.Alg) {
		return jwt_errors.ErrUnsupportedAlgorithm
	}

	return nil
}

func NewHeader(alg string) *Header {
	return &Header{
		Alg: alg,
		Typ: "JWT",
	}
}
