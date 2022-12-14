package servergrpc

import (
	"GophKeeper/internal/server/servergrpc/services/auth_service"
	"GophKeeper/internal/storage"
	pb "GophKeeper/pkg/proto/auth"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"time"
)

type ServerGRPC struct {
	*grpc.Server
	net.Listener
	pb.AuthServiceServer
}

func NewServer(addr string, store storage.UserStorage) (*ServerGRPC, error) {

	listen, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	s := &ServerGRPC{
		Server:   grpc.NewServer(),
		Listener: listen,
	}

	registerServices(s.Server, store)

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

func registerServices(server *grpc.Server, store storage.UserStorage) {

	authServ := auth_service.NewAuthService(store)

	pb.RegisterAuthServiceServer(server, authServ)
}
