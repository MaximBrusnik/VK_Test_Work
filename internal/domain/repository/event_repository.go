package repository

import (
	"context"

	"awesomeProject3/internal/domain/entity"
)

// EventRepository defines the interface for event storage
type EventRepository interface {
	// Save saves an event to the repository
	Save(ctx context.Context, event *entity.Event) error

	// FindByKey finds all events for a given key
	FindByKey(ctx context.Context, key string) ([]*entity.Event, error)

	// Subscribe subscribes to events for a given key
	Subscribe(ctx context.Context, key string, handler func(*entity.Event)) error

	// Unsubscribe unsubscribes from events for a given key
	Unsubscribe(ctx context.Context, key string) error

	// Close closes the repository
	Close(ctx context.Context) error
} 