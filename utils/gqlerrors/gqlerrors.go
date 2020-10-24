package gqlerrors

import (
	"strings"

	"github.com/pkg/errors"
)

var _ error = (*SQLError)(nil)

type SQLError struct {
	error
}

var (
	UniqueError = SQLError{error: errors.New("Unique constraint failed")}
)

// Проверка ошибки запроса к БД
func Is(err error, target SQLError) bool {
	return strings.Contains(err.Error(), target.Error())
}
