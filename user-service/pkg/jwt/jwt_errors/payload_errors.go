package jwt_errors

import "errors"

var (
	ErrInvalidPayload       = errors.New("неверная полезная нагрузка")
	ErrInvalidPayloadFormat = errors.New("неверный формат полезной нагрузки")
	ErrInvalidTimeRange     = errors.New("время истечения меньше времени начала действия")
	ErrTokenExpired         = errors.New("токен истек")
	ErrNotBeforeError       = errors.New("токен еще не активен (nbf)")
)
