package clientgrpc

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "GophKeeper/pkg/proto/auth"
)

type ClientGRPC struct {
	addr      string
	conn      *grpc.ClientConn
	rpcClient pb.AuthServiceClient
}

func NewClient(addr string) *ClientGRPC {
	return &ClientGRPC{
		addr: addr,
	}
}

func (c ClientGRPC) Register() error {
	_, err := c.rpcClient.Register(context.Background(), &pb.AuthRequest{
		Email:    "rnikbond@yandex.ru",
		Password: "qwerty123",
	})

	return err
}

func (c ClientGRPC) Login() error {
	_, err := c.rpcClient.Login(context.Background(), &pb.AuthRequest{
		Email:    "rnikbond@yandex.ru",
		Password: "qwerty123",
	})

	return err
}

func (c *ClientGRPC) Connect() (err error) {

	c.conn, err = grpc.Dial(c.addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	c.rpcClient = pb.NewAuthServiceClient(c.conn)
	return err
}

func (c *ClientGRPC) Disconnect() error {

	if c.conn == nil {
		return nil
	}

	return c.conn.Close()
}