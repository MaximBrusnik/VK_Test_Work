package subscribe

import (
	"awesomeProject3/internal/domain/entity"
	"awesomeProject3/internal/domain/errors"
	"awesomeProject3/internal/domain/repository"
	"context"
	"github.com/sirupsen/logrus"
)

// UseCase defines the subscription use case interface
type UseCase interface {
	Execute(ctx context.Context, req Request, callback func(*entity.Event)) error
}

// Request represents a subscription request
type Request struct {
	Key string
}

// subscribeUseCase implements the subscription use case
type subscribeUseCase struct {
	eventRepo repository.EventRepository
	logger    *logrus.Logger
}

// New creates a new subscription use case
func New(eventRepo repository.EventRepository, logger *logrus.Logger) UseCase {
	return &subscribeUseCase{
		eventRepo: eventRepo,
		logger:    logger,
	}
}

// Execute handles the subscription request
func (uc *subscribeUseCase) Execute(ctx context.Context, req Request, callback func(*entity.Event)) error {
	if req.Key == "" {
		return errors.ErrInvalidEventKey
	}

	uc.logger.WithField("key", req.Key).Info("subscribing to events")

	// Subscribe to events
	err := uc.eventRepo.Subscribe(ctx, req.Key, callback)
	if err != nil {
		uc.logger.WithError(err).Error("failed to subscribe to events")
		return err
	}

	return nil
}
