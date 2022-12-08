package servergrpc

import (
	"fmt"
	"google.golang.org/grpc"
	"net"
	"time"
)

type ServerGRPC struct {
	*grpc.Server
	net.Listener
}

func NewServer(addr string) (*ServerGRPC, error) {

	listen, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &ServerGRPC{
		Server:   grpc.NewServer(),
		Listener: listen,
	}, nil
}

func (s *ServerGRPC) Start() {

	go func() {
		fmt.Printf("Server started at: %s\n", time.Now().Format("02-01-2006 15:04:05"))

		if err := s.Server.Serve(s.Listener); err != nil {
			fmt.Printf("failed start gRPC server: %v\n", err)
		}
	}()
}

func (s *ServerGRPC) Stop() {

	s.Server.Stop()
	fmt.Printf("Server stopped at: %s\n", time.Now().Format("02-01-2006 15:04:05"))
}
