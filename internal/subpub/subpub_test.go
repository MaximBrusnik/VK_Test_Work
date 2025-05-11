package subpub

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestNewSubPub(t *testing.T) {
	sp := NewSubPub()
	if sp == nil {
		t.Error("NewSubPub returned nil")
	}
}

func TestSubscribeAndPublish(t *testing.T) {
	sp := NewSubPub()
	received := make(chan interface{}, 1)

	sub, err := sp.Subscribe("test", func(msg interface{}) {
		received <- msg
	})
	if err != nil {
		t.Fatalf("Subscribe failed: %v", err)
	}

	msg := "test message"
	err = sp.Publish("test", msg)
	if err != nil {
		t.Fatalf("Publish failed: %v", err)
	}

	select {
	case got := <-received:
		if got != msg {
			t.Errorf("got %v, want %v", got, msg)
		}
	case <-time.After(time.Second):
		t.Error("timeout waiting for message")
	}

	sub.Unsubscribe()
}

func TestMultipleSubscribers(t *testing.T) {
	sp := NewSubPub()
	var wg sync.WaitGroup
	subscriberCount := 3
	received := make([]chan interface{}, subscriberCount)

	for i := 0; i < subscriberCount; i++ {
		received[i] = make(chan interface{}, 1)
		wg.Add(1)
		idx := i
		sub, err := sp.Subscribe("test", func(msg interface{}) {
			received[idx] <- msg
			wg.Done()
		})
		if err != nil {
			t.Fatalf("Subscribe failed: %v", err)
		}
		defer sub.Unsubscribe()
	}

	msg := "test message"
	err := sp.Publish("test", msg)
	if err != nil {
		t.Fatalf("Publish failed: %v", err)
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// All subscribers received the message
		for i := 0; i < subscriberCount; i++ {
			select {
			case got := <-received[i]:
				if got != msg {
					t.Errorf("subscriber %d got %v, want %v", i, got, msg)
				}
			case <-time.After(time.Second):
				t.Errorf("timeout waiting for subscriber %d", i)
			}
		}
	case <-time.After(time.Second):
		t.Error("timeout waiting for all subscribers")
	}
}

func TestUnsubscribe(t *testing.T) {
	sp := NewSubPub()
	received := make(chan interface{}, 1)

	sub, err := sp.Subscribe("test", func(msg interface{}) {
		received <- msg
	})
	if err != nil {
		t.Fatalf("Subscribe failed: %v", err)
	}

	sub.Unsubscribe()

	err = sp.Publish("test", "test message")
	if err != nil {
		t.Fatalf("Publish failed: %v", err)
	}

	select {
	case <-received:
		t.Error("received message after unsubscribe")
	case <-time.After(100 * time.Millisecond):
		// Expected timeout - no message should be received
	}
}

func TestClose(t *testing.T) {
	sp := NewSubPub()

	// Создаем подписку перед закрытием
	received := make(chan interface{}, 1)
	_, err := sp.Subscribe("test", func(msg interface{}) {
		received <- msg
	})
	if err != nil {
		t.Fatalf("Subscribe failed: %v", err)
	}

	// Проверяем, что подписка работает до закрытия
	err = sp.Publish("test", "test message")
	if err != nil {
		t.Errorf("Publish before close failed: %v", err)
	}

	select {
	case msg := <-received:
		if msg != "test message" {
			t.Errorf("received wrong message: got %v, want %v", msg, "test message")
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("timeout waiting for message before close")
	}

	// Создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// Закрываем сервис
	err = sp.Close(ctx)
	if err != context.DeadlineExceeded {
		t.Errorf("Close error: got %v, want %v", err, context.DeadlineExceeded)
	}

	// Проверяем повторное закрытие
	err = sp.Close(ctx)
	if err != nil {
		t.Errorf("Second close failed: %v", err)
	}

	// Проверяем, что подписка больше не работает
	err = sp.Publish("test", "test message")
	if err != ErrClosed {
		t.Errorf("Publish after close: got %v, want %v", err, ErrClosed)
	}

	select {
	case <-received:
		t.Error("received message after close")
	case <-time.After(100 * time.Millisecond):
		// Ожидаемый таймаут - сообщение не должно быть получено
	}

	// Пробуем создать новую подписку после закрытия
	_, err = sp.Subscribe("test", func(msg interface{}) {})
	if err != ErrClosed {
		t.Errorf("Subscribe after close: got %v, want %v", err, ErrClosed)
	}
}

func TestSlowSubscriber(t *testing.T) {
	sp := NewSubPub()
	var wg sync.WaitGroup
	fastReceived := make(chan interface{}, 1)
	slowReceived := make(chan interface{}, 1)

	// Fast subscriber
	wg.Add(1)
	fastSub, err := sp.Subscribe("test", func(msg interface{}) {
		fastReceived <- msg
		wg.Done()
	})
	if err != nil {
		t.Fatalf("Subscribe failed: %v", err)
	}
	defer fastSub.Unsubscribe()

	// Slow subscriber
	wg.Add(1)
	slowSub, err := sp.Subscribe("test", func(msg interface{}) {
		time.Sleep(500 * time.Millisecond)
		slowReceived <- msg
		wg.Done()
	})
	if err != nil {
		t.Fatalf("Subscribe failed: %v", err)
	}
	defer slowSub.Unsubscribe()

	msg := "test message"
	err = sp.Publish("test", msg)
	if err != nil {
		t.Fatalf("Publish failed: %v", err)
	}

	// Fast subscriber should receive message quickly
	select {
	case got := <-fastReceived:
		if got != msg {
			t.Errorf("fast subscriber got %v, want %v", got, msg)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("timeout waiting for fast subscriber")
	}

	// Slow subscriber should still receive message
	select {
	case got := <-slowReceived:
		if got != msg {
			t.Errorf("slow subscriber got %v, want %v", got, msg)
		}
	case <-time.After(time.Second):
		t.Error("timeout waiting for slow subscriber")
	}

	wg.Wait()
}
