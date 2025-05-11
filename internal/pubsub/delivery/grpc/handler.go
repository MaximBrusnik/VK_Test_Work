package grpc

import (
	"context"
	"sync"

	"awesomeProject3/internal/domain/entity"
	"awesomeProject3/internal/usecase/publish"
	"awesomeProject3/internal/usecase/subscribe"
	"awesomeProject3/pkg/proto"
	"awesomeProject3/pkg/validator"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Handler реализует gRPC сервер для PubSub
type Handler struct {
	proto.UnimplementedPubSubServer
	logger      *logrus.Logger
	publishUC   publish.UseCase
	subscribeUC subscribe.UseCase
	mu          sync.RWMutex
	subs        map[string]map[chan *proto.Event]struct{}
}

// NewHandler создает новый обработчик gRPC
func NewHandler(
	logger *logrus.Logger,
	publishUC publish.UseCase,
	subscribeUC subscribe.UseCase,
) *Handler {
	return &Handler{
		logger:      logger,
		publishUC:   publishUC,
		subscribeUC: subscribeUC,
		subs:        make(map[string]map[chan *proto.Event]struct{}),
	}
}

// Subscribe обрабатывает запрос на подписку
func (h *Handler) Subscribe(req *proto.SubscribeRequest, stream proto.PubSub_SubscribeServer) error {
	if err := validator.ValidateNotEmpty(req.GetKey(), "key"); err != nil {
		return status.Error(codes.InvalidArgument, "требуется указать ключ")
	}

	key := req.GetKey()

	// Создаем канал для этой подписки
	events := make(chan *proto.Event, 100)
	defer close(events)

	// Регистрируем подписку
	h.mu.Lock()
	if _, exists := h.subs[key]; !exists {
		h.subs[key] = make(map[chan *proto.Event]struct{})
	}
	h.subs[key][events] = struct{}{}
	h.mu.Unlock()

	// Очистка при выходе
	defer func() {
		h.mu.Lock()
		delete(h.subs[key], events)
		if len(h.subs[key]) == 0 {
			delete(h.subs, key)
		}
		h.mu.Unlock()
	}()

	// Подписываемся используя use case
	err := h.subscribeUC.Execute(stream.Context(), subscribe.Request{Key: key}, func(event *entity.Event) {
		protoEvent := &proto.Event{Data: event.Data}
		select {
		case events <- protoEvent:
		default:
			h.logger.WithField("key", key).Warn("буфер подписчика переполнен, сообщение отброшено")
		}
	})
	if err != nil {
		return status.Error(codes.Internal, "не удалось подписаться")
	}

	// Отправляем события клиенту
	for {
		select {
		case event, ok := <-events:
			if !ok {
				return nil
			}
			if err := stream.Send(event); err != nil {
				return status.Error(codes.Internal, "не удалось отправить событие")
			}
		case <-stream.Context().Done():
			return nil
		}
	}
}

// Publish обрабатывает запрос на публикацию
func (h *Handler) Publish(ctx context.Context, req *proto.PublishRequest) (*emptypb.Empty, error) {
	if err := validator.ValidateAll(
		func() error { return validator.ValidateNotEmpty(req.GetKey(), "key") },
		func() error { return validator.ValidateNotEmpty(req.GetData(), "data") },
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Публикуем используя use case
	err := h.publishUC.Execute(ctx, publish.Request{
		Key:  req.GetKey(),
		Data: req.GetData(),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "не удалось опубликовать")
	}

	return &emptypb.Empty{}, nil
}
