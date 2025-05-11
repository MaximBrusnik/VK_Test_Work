package publish

import (
	"context"

	"awesomeProject3/internal/domain/entity"
	"awesomeProject3/internal/domain/repository"
	"github.com/sirupsen/logrus"
)

// Request represents the publish request
type Request struct {
	Key  string
	Data string
}

// UseCase defines the publish use case
type UseCase interface {
	Execute(ctx context.Context, req Request) error
}

type publishUseCase struct {
	eventRepo repository.EventRepository
	logger    *logrus.Logger
}

// New creates a new publish use case
func New(eventRepo repository.EventRepository, logger *logrus.Logger) UseCase {
	return &publishUseCase{
		eventRepo: eventRepo,
		logger:    logger,
	}
}

// Execute executes the publish use case
func (uc *publishUseCase) Execute(ctx context.Context, req Request) error {
	// Create and validate event
	event := entity.NewEvent(req.Key, req.Data)
	if err := event.Validate(); err != nil {
		uc.logger.WithError(err).WithFields(logrus.Fields{
			"key":  req.Key,
			"data": req.Data,
		}).Error("failed to validate event")
		return err
	}

	// Save event
	if err := uc.eventRepo.Save(ctx, event); err != nil {
		uc.logger.WithError(err).WithFields(logrus.Fields{
			"event_id": event.ID,
			"key":      event.Key,
		}).Error("failed to save event")
		return err
	}

	uc.logger.WithFields(logrus.Fields{
		"event_id": event.ID,
		"key":      event.Key,
	}).Info("event published successfully")

	return nil
} 