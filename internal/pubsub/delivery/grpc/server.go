package grpc

import (
	"fmt"
	"net"

	"awesomeProject3/pkg/grpc/middleware"
	"awesomeProject3/pkg/proto"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// Server представляет gRPC сервер
type Server struct {
	server *grpc.Server
	logger *logrus.Logger
	port   int
}

// NewServer создает новый gRPC сервер
func NewServer(handler *Handler, logger *logrus.Logger, port int) *Server {
	server := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.LoggingInterceptor(logger)),
		grpc.StreamInterceptor(middleware.StreamLoggingInterceptor(logger)),
	)
	proto.RegisterPubSubServer(server, handler)

	return &Server{
		server: server,
		logger: logger,
		port:   port,
	}
}

// Start запускает gRPC сервер
func (s *Server) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return fmt.Errorf("не удалось запустить сервер: %v", err)
	}

	s.logger.WithField("port", s.port).Info("запуск gRPC сервера")

	if err := s.server.Serve(lis); err != nil {
		return fmt.Errorf("ошибка при работе сервера: %v", err)
	}

	return nil
}

// Stop останавливает gRPC сервер
func (s *Server) Stop() {
	s.logger.Info("остановка gRPC сервера")
	s.server.GracefulStop()
} 