package servergrpc

import (
	"fmt"
	"net"
	"time"

	"google.golang.org/grpc"

	"GophKeeper/internal/server/model"
	"GophKeeper/internal/server/servergrpc/services/auth_service"
	pb "GophKeeper/pkg/proto/auth"
)

type ServerGRPC struct {
	*grpc.Server
	net.Listener

	auth *auth_service.AuthService
}

func NewServer(addr string, auth *model.AuthModel) (*ServerGRPC, error) {

	listen, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	s := &ServerGRPC{
		Server:   grpc.NewServer(),
		Listener: listen,
		auth:     auth_service.NewAuthService(auth),
	}

	pb.RegisterAuthServiceServer(s.Server, s.auth)

	return s, nil
}

func (serv *ServerGRPC) Start() {

	go func() {
		fmt.Printf("Server started at: %s\n", time.Now().Format("02-01-2006 15:04:05"))

		if err := serv.Server.Serve(serv.Listener); err != nil {
			fmt.Printf("failed start gRPC server: %v\n", err)
		}
	}()
}

func (serv *ServerGRPC) Stop() {

	serv.Server.Stop()
	fmt.Printf("Server stopped at: %s\n", time.Now().Format("02-01-2006 15:04:05"))
}
