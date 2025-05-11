package entity

import (
	"awesomeProject3/internal/domain/errors"
	"awesomeProject3/pkg/validator"
	"time"
)

// Event представляет собой доменное событие в системе
type Event struct {
	ID        string
	Key       string
	Data      string
	Timestamp time.Time
}

// NewEvent создает новое событие
func NewEvent(key, data string) *Event {
	return &Event{
		ID:        generateID(),
		Key:       key,
		Data:      data,
		Timestamp: time.Now(),
	}
}

// Validate проверяет корректность события
func (e *Event) Validate() error {
	return validator.ValidateAll(
		func() error {
			if err := validator.ValidateNotEmpty(e.Key, "key"); err != nil {
				return errors.ErrInvalidEventKey
			}
			return nil
		},
		func() error {
			if err := validator.ValidateNotEmpty(e.Data, "data"); err != nil {
				return errors.ErrInvalidEventData
			}
			return nil
		},
		func() error {
			if err := validator.ValidateLength(e.Key, 1, 100); err != nil {
				return errors.ErrInvalidEventKey
			}
			return nil
		},
	)
}

// generateID генерирует уникальный идентификатор для события
func generateID() string {
	return time.Now().Format("20060102150405.000000000")
}
