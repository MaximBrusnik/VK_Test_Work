package subpub

import (
	"context"
	"sync"
)

type MessageHandler func(msg interface{})

type Subscription interface {
	Unsubscribe()
}

type subscription struct {
	subject string
	handler MessageHandler
	done    chan struct{}
}

func (s *subscription) Unsubscribe() {
	close(s.done)
}

type SubPub interface {
	Subscribe(subject string, cb MessageHandler) (Subscription, error)
	Publish(subject string, msg interface{}) error
	Close(ctx context.Context) error
}

type subPub struct {
	mu          sync.RWMutex
	subscribers map[string][]*subscription
	closed      bool
}

func NewSubPub() SubPub {
	return &subPub{
		subscribers: make(map[string][]*subscription),
	}
}

func (sp *subPub) Subscribe(subject string, cb MessageHandler) (Subscription, error) {
	sp.mu.Lock()
	defer sp.mu.Unlock()

	if sp.closed {
		return nil, ErrClosed
	}

	sub := &subscription{
		subject: subject,
		handler: cb,
		done:    make(chan struct{}),
	}

	sp.subscribers[subject] = append(sp.subscribers[subject], sub)
	return sub, nil
}

func (sp *subPub) Publish(subject string, msg interface{}) error {
	sp.mu.RLock()
	if sp.closed {
		sp.mu.RUnlock()
		return ErrClosed
	}

	subscribers := make([]*subscription, len(sp.subscribers[subject]))
	copy(subscribers, sp.subscribers[subject])
	sp.mu.RUnlock()

	for _, sub := range subscribers {
		select {
		case <-sub.done:
			continue
		default:
			go func(s *subscription) {
				select {
				case <-s.done:
					return
				default:
					s.handler(msg)
				}
			}(sub)
		}
	}

	return nil
}

func (sp *subPub) Close(ctx context.Context) error {
	sp.mu.Lock()
	if sp.closed {
		sp.mu.Unlock()
		return nil
	}
	sp.closed = true

	// Очищаем все подписки
	for _, subs := range sp.subscribers {
		for _, sub := range subs {
			select {
			case <-sub.done:
				// Канал уже закрыт, пропускаем
				continue
			default:
				close(sub.done)
			}
		}
	}
	sp.subscribers = make(map[string][]*subscription)
	sp.mu.Unlock()

	// Ждем отмены контекста или таймаута
	select {
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Custom errors
var (
	ErrClosed = &Error{"subpub: service is closed"}
)

type Error struct {
	msg string
}

func (e *Error) Error() string {
	return e.msg
}
