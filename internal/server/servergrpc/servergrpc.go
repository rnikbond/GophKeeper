package servergrpc

import (
	"net"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"GophKeeper/internal/server/servergrpc/rpc_services"
	pbAuth "GophKeeper/pkg/proto/auth"
	pbCred "GophKeeper/pkg/proto/data/credential"
)

// ServerOption - определяет операцию сервиса авторизации.
type ServerOption func(serv *ServerGRPC)

// ServerGRPC структура gPRC сервера.
type ServerGRPC struct {
	*grpc.Server
	net.Listener

	auth   *rpc_services.AuthServiceRPC
	logger *zap.Logger

	secretKey string
}

// NewServer - Создание экземпляра gRPC сервера, но не запускает его.
// • addr - Адрес, на котором в при вызове Start() будет запущен сервер.
func NewServer(addr string, interceptor grpc.ServerOption, opts ...ServerOption) (*ServerGRPC, error) {
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	s := &ServerGRPC{
		//Server:   grpc.NewServer(interceptors.NewValidateInterceptor(secretKey)),
		Server:   grpc.NewServer(interceptor),
		Listener: listen,
		logger:   zap.L(),
	}

	for _, opt := range opts {
		opt(s)
	}

	return s, nil
}

// WithAuthServiceRPC - Регистрирует сервис gPRC авторизации
func WithAuthServiceRPC(auth *rpc_services.AuthServiceRPC) ServerOption {
	return func(serv *ServerGRPC) {
		pbAuth.RegisterAuthServiceServer(serv.Server, auth)
	}
}

// WithCredServiceRPC - Регистрирует сервис gPRC для хранения логинов и паролей
func WithCredServiceRPC(cred *rpc_services.CredServiceRPC) ServerOption {
	return func(serv *ServerGRPC) {
		pbCred.RegisterCredentialServiceServer(serv.Server, cred)
	}
}

// Start - Запуск сервера.
func (serv *ServerGRPC) Start() {
	go func() {
		serv.logger.Info("Server started", zap.String("At", time.Now().Format("02-01-2006 15:04:05")))
		if err := serv.Server.Serve(serv.Listener); err != nil {
			serv.logger.Error("failed run gRPC server", zap.Error(err))
		}
	}()
}

// Stop - Остановка сервера.
func (serv *ServerGRPC) Stop() {
	serv.Server.Stop()
	serv.logger.Info("Server stopped", zap.String("At", time.Now().Format("02-01-2006 15:04:05")))
}
