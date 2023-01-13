package server_grpc

import (
	"GophKeeper/internal/server/server_grpc/services/grpc_service_auth"
	"GophKeeper/internal/server/server_grpc/services/grpc_service_binary"
	"GophKeeper/internal/server/server_grpc/services/grpc_service_card"
	"GophKeeper/internal/server/server_grpc/services/grpc_service_cred"
	"GophKeeper/internal/server/server_grpc/services/grpc_service_text"
	pbBinary "GophKeeper/pkg/proto/binary"
	pbCard "GophKeeper/pkg/proto/card"
	pbCred "GophKeeper/pkg/proto/credential"
	pbText "GophKeeper/pkg/proto/text"
	"net"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	pbAuth "GophKeeper/pkg/proto/auth"
)

// ServerOption - определяет операцию сервиса авторизации.
type ServerOption func(serv *ServerGRPC)

// ServerGRPC структура gPRC сервера.
type ServerGRPC struct {
	*grpc.Server
	net.Listener

	auth   *grpc_service_auth.AuthServiceRPC
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
func WithAuthServiceRPC(auth *grpc_service_auth.AuthServiceRPC) ServerOption {
	return func(serv *ServerGRPC) {
		pbAuth.RegisterAuthServiceServer(serv.Server, auth)
	}
}

// WithCredServiceRPC - Регистрирует сервис gPRC для хранения логинов и паролей
func WithCredServiceRPC(cred *grpc_service_cred.CredServiceRPC) ServerOption {
	return func(serv *ServerGRPC) {
		pbCred.RegisterCredentialServiceServer(serv.Server, cred)
	}
}

// WithBinaryServiceRPC - Регистрирует сервис gPRC для хранения бинарных данных
func WithBinaryServiceRPC(bin *grpc_service_binary.BinaryServiceRPC) ServerOption {
	return func(serv *ServerGRPC) {
		pbBinary.RegisterBinaryServiceServer(serv.Server, bin)
	}
}

// WithTextServiceRPC - Регистрирует сервис gPRC для хранения текстовых данных
func WithTextServiceRPC(txt *grpc_service_text.TextServiceRPC) ServerOption {
	return func(serv *ServerGRPC) {
		pbText.RegisterTextServiceServer(serv.Server, txt)
	}
}

func WithCardServiceRPC(cardServ *grpc_service_card.CardServiceRPC) ServerOption {
	return func(serv *ServerGRPC) {
		pbCard.RegisterCardServiceServer(serv.Server, cardServ)
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
