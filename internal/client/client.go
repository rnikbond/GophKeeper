package client

type Client interface {
	Connect(addr string) error
	Disconnect() error
}
