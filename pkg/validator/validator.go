package validator

import (
	"errors"
	"strings"
)

var (
	// ErrEmptyValue возвращается, когда обязательное значение пустое
	ErrEmptyValue = errors.New("значение не может быть пустым")

	// ErrInvalidLength возвращается, когда длина значения некорректна
	ErrInvalidLength = errors.New("некорректная длина")
)

// ValidateNotEmpty проверяет, что строка не пустая
func ValidateNotEmpty(value, fieldName string) error {
	if strings.TrimSpace(value) == "" {
		return ErrEmptyValue
	}
	return nil
}

// ValidateLength проверяет, что длина строки находится в указанном диапазоне
func ValidateLength(value string, min, max int) error {
	length := len(strings.TrimSpace(value))
	if length < min || length > max {
		return ErrInvalidLength
	}
	return nil
}

// ValidateAll запускает все валидаторы и возвращает первую ошибку
func ValidateAll(validators ...func() error) error {
	for _, validator := range validators {
		if err := validator(); err != nil {
			return err
		}
	}
	return nil
} 