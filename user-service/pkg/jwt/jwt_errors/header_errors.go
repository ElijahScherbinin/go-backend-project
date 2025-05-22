package jwt_errors

import "errors"

var (
	ErrInvalidHeader        = errors.New("неверный заголовок")
	ErrInvalidHeaderFormat  = errors.New("неверный формат заголовка")
	ErrUnsupportedTypeToken = errors.New("неподдерживаемый тип токена")
	ErrMissingAlgorithm     = errors.New("алгоритм подписи не указан")
	ErrUnsupportedAlgorithm = errors.New("неподдерживаемый алгоритм подписи")
)
