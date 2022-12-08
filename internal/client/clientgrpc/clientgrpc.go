package clientgrpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ClientGRPC struct {
	addr string
	conn *grpc.ClientConn
}

func NewClient(addr string) *ClientGRPC {
	return &ClientGRPC{
		addr: addr,
	}
}

func (c *ClientGRPC) Connect() (err error) {

	c.conn, err = grpc.Dial(c.addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	return err
}

func (c *ClientGRPC) Disconnect() error {

	if c.conn == nil {
		return nil
	}

	return c.conn.Close()
}
