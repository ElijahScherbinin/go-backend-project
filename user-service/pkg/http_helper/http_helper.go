package http_helper

import (
	"errors"
	"net/url"
	"strconv"
)

var ErrParamIsEmpty error = errors.New("query param is empty")
var ErrUnsupportedType error = errors.New("unsupported type")

func GetQueryParam[T any](values url.Values, key string) (T, error) {
	param := values.Get(key)
	if param == "" {
		return *new(T), ErrParamIsEmpty
	}
	return convertParam[T](param)
}

func GetRouteParam[T any](params map[string]string, key string) (T, error) {
	param, isExist := params[key]
	if param == "" || !isExist {
		return *new(T), ErrParamIsEmpty
	}
	return convertParam[T](param)
}

func convertParam[T any](param string) (T, error) {
	var result = new(T)
	var err error = nil

	switch value := any(result).(type) {
	case *int:
		*value, err = strconv.Atoi(param)
	case *string:
		*value = param
	default:
		err = ErrUnsupportedType
	}

	return *result, err
}
