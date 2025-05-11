package repository

import (
	"context"
	"sync"

	"awesomeProject3/internal/domain/entity"
	"awesomeProject3/internal/domain/errors"
)

// InMemoryRepository implements EventRepository interface using in-memory storage
type InMemoryRepository struct {
	mu          sync.RWMutex
	events      map[string][]*entity.Event
	subscribers map[string][]func(*entity.Event)
	closed      bool
}

// NewInMemoryRepository creates a new in-memory repository
func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		events:      make(map[string][]*entity.Event),
		subscribers: make(map[string][]func(*entity.Event)),
	}
}

// Save saves an event to the repository
func (r *InMemoryRepository) Save(ctx context.Context, event *entity.Event) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.closed {
		return errors.ErrServiceClosed
	}

	// Save event
	r.events[event.Key] = append(r.events[event.Key], event)

	// Notify subscribers
	if subscribers, ok := r.subscribers[event.Key]; ok {
		for _, handler := range subscribers {
			handler(event)
		}
	}

	return nil
}

// FindByKey finds all events for a given key
func (r *InMemoryRepository) FindByKey(ctx context.Context, key string) ([]*entity.Event, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.closed {
		return nil, errors.ErrServiceClosed
	}

	events, ok := r.events[key]
	if !ok {
		return nil, errors.ErrEventNotFound
	}

	return events, nil
}

// Subscribe subscribes to events for a given key
func (r *InMemoryRepository) Subscribe(ctx context.Context, key string, handler func(*entity.Event)) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.closed {
		return errors.ErrServiceClosed
	}

	if key == "" {
		return errors.ErrInvalidEventKey
	}

	r.subscribers[key] = append(r.subscribers[key], handler)
	return nil
}

// Unsubscribe unsubscribes from events for a given key
func (r *InMemoryRepository) Unsubscribe(ctx context.Context, key string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.closed {
		return errors.ErrServiceClosed
	}

	delete(r.subscribers, key)
	return nil
}

// Close closes the repository
func (r *InMemoryRepository) Close(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.closed {
		return nil
	}

	r.closed = true
	return nil
}
