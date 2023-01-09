package clientgrpc

import (
	"GophKeeper/internal/model/binary"
	"GophKeeper/internal/model/card"
	"GophKeeper/internal/model/text"
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"GophKeeper/internal/model/cred"
	pbAuth "GophKeeper/pkg/proto/auth"
	pbBinary "GophKeeper/pkg/proto/data/binary"
	pbCard "GophKeeper/pkg/proto/data/card"
	pbCred "GophKeeper/pkg/proto/data/credential"
	pbText "GophKeeper/pkg/proto/data/text"
)

type ClientGRPC struct {
	addr            string
	token           string
	conn            *grpc.ClientConn
	rpcAuthClient   pbAuth.AuthServiceClient
	rpcCredClient   pbCred.CredentialServiceClient
	rpcBinaryClient pbBinary.BinaryServiceClient
	rpcTextClient   pbText.TextServiceClient
	rpcCardClient   pbCard.CardServiceClient
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

func (c ClientGRPC) CreateBinary() error {

	md := metadata.New(map[string]string{"token": c.token})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	_, err := c.rpcBinaryClient.Create(ctx, &pbBinary.CreateRequest{
		MetaInfo: "www.ololo.com",
		Data:     []byte("123123123123asd"),
	})

	return err
}

func (c ClientGRPC) FindBinary() (binary.DataFull, error) {

	md := metadata.New(map[string]string{"token": c.token})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	out, err := c.rpcBinaryClient.Get(ctx, &pbBinary.GetRequest{
		MetaInfo: "www.ololo.com",
	})

	data := binary.DataFull{
		MetaInfo: out.MetaInfo,
		Bytes:    out.Data,
	}

	return data, err
}

func (c ClientGRPC) CreateText() error {

	md := metadata.New(map[string]string{"token": c.token})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	_, err := c.rpcTextClient.Create(ctx, &pbText.CreateRequest{
		MetaInfo: "book",
		Text:     "123123123123asd",
	})

	return err
}

func (c ClientGRPC) FindText() (text.DataTextFull, error) {

	md := metadata.New(map[string]string{"token": c.token})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	out, err := c.rpcTextClient.Get(ctx, &pbText.GetRequest{
		MetaInfo: "book",
	})

	data := text.DataTextFull{
		MetaInfo: out.MetaInfo,
		Text:     out.Text,
	}

	return data, err
}

func (c ClientGRPC) CreateCard() error {

	md := metadata.New(map[string]string{"token": c.token})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	_, err := c.rpcCardClient.Create(ctx, &pbCard.CreateRequest{
		MetaInfo: "MirPay",
		Number:   "4648289760410976",
		Period:   "10.2030",
		CVV:      "111",
		FullName: "Test Test",
	})

	return err
}

func (c ClientGRPC) FindCard() (card.DataCard, error) {

	md := metadata.New(map[string]string{"token": c.token})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	out, err := c.rpcCardClient.Get(ctx, &pbCard.GetRequest{
		MetaInfo: "MirPay",
	})

	data := card.DataCard{
		MetaInfo: "MirPay",
		Number:   out.Number,
		Period:   out.Period,
		CVV:      out.CVV,
		FullName: out.FullName,
	}

	return data, err
}

func (c *ClientGRPC) Connect() (err error) {

	c.conn, err = grpc.Dial(c.addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	c.rpcAuthClient = pbAuth.NewAuthServiceClient(c.conn)
	c.rpcCredClient = pbCred.NewCredentialServiceClient(c.conn)
	c.rpcBinaryClient = pbBinary.NewBinaryServiceClient(c.conn)
	c.rpcTextClient = pbText.NewTextServiceClient(c.conn)
	c.rpcCardClient = pbCard.NewCardServiceClient(c.conn)
	return err
}

func (c *ClientGRPC) Disconnect() error {

	if c.conn == nil {
		return nil
	}

	return c.conn.Close()
}
