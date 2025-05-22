package jwt_errors

import "errors"

var (
	ErrTokenGeneration             = errors.New("ошибка генерации токена")
	ErrConvertToJson               = errors.New("ошибка конветации в JSON")
	ErrInvalidTokenFormat          = errors.New("неверный формат токена")
	ErrExtractTokenIsEmpty         = errors.New("извлеченный токен пустой")
	ErrSignatureVerificationFailed = errors.New("проверка подписи не удалась")
)
