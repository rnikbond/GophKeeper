package clientgrpc

import (
	"GophKeeper/internal/model/cred"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	pbAuth "GophKeeper/pkg/proto/auth"
	pbCred "GophKeeper/pkg/proto/data/credential"
)

type ClientGRPC struct {
	addr          string
	token         string
	conn          *grpc.ClientConn
	rpcAuthClient pbAuth.AuthServiceClient
	rpcCredClient pbCred.CredentialServiceClient
}

func NewClient(addr string) *ClientGRPC {
	return &ClientGRPC{
		addr: addr,
	}
}

func (c *ClientGRPC) Login() error {
	resp, err := c.rpcAuthClient.Login(context.Background(), &pbAuth.AuthRequest{
		Email:    "rnikbond@yandex.ru",
		Password: "qwerty123",
	})

	if resp != nil {
		c.token = resp.Token
		fmt.Printf("token: %s\n", resp.Token)
	}

	return err
}

func (c *ClientGRPC) Register() error {
	resp, err := c.rpcAuthClient.Register(context.Background(), &pbAuth.AuthRequest{
		Email:    "rnikbond@yandex.ru",
		Password: "qwerty123",
	})

	if resp != nil {
		c.token = resp.Token
		fmt.Printf("token: %s\n", resp.Token)
	}

	return err
}

func (c ClientGRPC) ChangePassword() error {

	md := metadata.New(map[string]string{"token": c.token})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	_, err := c.rpcAuthClient.ChangePassword(ctx, &pbAuth.ChangePasswordRequest{
		Password: "qwerty123",
	})

	return err
}

func (c ClientGRPC) CreatePairCred() error {

	md := metadata.New(map[string]string{"token": c.token})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	_, err := c.rpcCredClient.Create(ctx, &pbCred.CreateRequest{
		Email:    "ololoev@email.com",
		Password: "qwerty123",
		MetaInfo: "www.ololo.com",
	})

	return err
}

func (c ClientGRPC) FindPairCred() (cred.CredentialFull, error) {

	md := metadata.New(map[string]string{"token": c.token})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	out, err := c.rpcCredClient.Get(ctx, &pbCred.GetRequest{
		Email:    "ololoev@email.com",
		MetaInfo: "www.ololo.com",
	})

	data := cred.CredentialFull{
		Email:    out.Email,
		MetaInfo: out.MetaInfo,
		Password: out.Password,
	}

	return data, err
}

func (c *ClientGRPC) Connect() (err error) {

	c.conn, err = grpc.Dial(c.addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	c.rpcAuthClient = pbAuth.NewAuthServiceClient(c.conn)
	c.rpcCredClient = pbCred.NewCredentialServiceClient(c.conn)
	return err
}

func (c *ClientGRPC) Disconnect() error {

	if c.conn == nil {
		return nil
	}

	return c.conn.Close()
}
