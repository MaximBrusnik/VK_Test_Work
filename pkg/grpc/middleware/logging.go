package middleware

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// LoggingInterceptor возвращает новый перехватчик для логирования unary запросов
func LoggingInterceptor(logger *logrus.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()

		// Вызываем обработчик
		resp, err := handler(ctx, req)

		// Логируем запрос
		duration := time.Since(start)
		entry := logger.WithFields(logrus.Fields{
			"method":   info.FullMethod,
			"duration": duration,
		})

		if err != nil {
			entry.WithError(err).Error("запрос завершился с ошибкой")
		} else {
			entry.Info("запрос успешно завершен")
		}

		return resp, err
	}
}

// StreamLoggingInterceptor возвращает новый перехватчик для логирования stream запросов
func StreamLoggingInterceptor(logger *logrus.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		start := time.Now()

		// Вызываем обработчик
		err := handler(srv, ss)

		// Логируем запрос
		duration := time.Since(start)
		entry := logger.WithFields(logrus.Fields{
			"method":   info.FullMethod,
			"duration": duration,
		})

		if err != nil {
			entry.WithError(err).Error("stream запрос завершился с ошибкой")
		} else {
			entry.Info("stream запрос успешно завершен")
		}

		return err
	}
} 