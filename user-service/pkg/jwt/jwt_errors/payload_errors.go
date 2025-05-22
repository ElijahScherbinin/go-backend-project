package jwt_errors

import "errors"

var (
	ErrInvalidClaims       = errors.New("неверная полезная нагрузка")
	ErrInvalidClaimsFormat = errors.New("неверный формат полезной нагрузки")
	ErrInvalidTimeRange    = errors.New("время истечения меньше времени начала действия")
	ErrTokenExpired        = errors.New("токен истек")
	ErrNotBeforeError      = errors.New("токен еще не активен (nbf)")
)
