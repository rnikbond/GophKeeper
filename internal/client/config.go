package client

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"strconv"
	"strings"
)

type Config struct {
	AddrGRPC   string `env:"ADDRESS_RPC" json:"address_rpc"`
	Salt       string `env:"SALT" json:"salt"`
	PublicKey  []byte `env:"PUBLIC_KEY" json:"public_key"`
	PrivateKey []byte `env:"PRIVATE_KEY" json:"private_key"`
}

// NewConfig Конфигурация сервера
func NewConfig() *Config {

	return &Config{
		AddrGRPC: ":3200",
		Salt:     "01.01.1970",
	}
}

// ParseArgs Разброр аргументов командной строки
func (cfg *Config) ParseArgs() error {

	addr := flag.String("ad", "", "string - address grpc gate")
	salt := flag.String("s", "", "string - password salt")
	privatePath := flag.String("prk", "", "private key - path to file")
	publicPath := flag.String("pbk", "", "public key - path to file")

	flag.Parse()

	if addr == nil || len(*addr) == 0 {
		*addr = cfg.AddrGRPC
	}

	if salt != nil && len(*salt) > 0 {
		cfg.Salt = *salt
	}

	if err := readKey(publicPath, &cfg.PublicKey); err != nil {
		return err
	}

	if err := readKey(privatePath, &cfg.PrivateKey); err != nil {
		return err
	}

	if err := isValidAddress(*addr); err != nil {
		return err
	}

	cfg.AddrGRPC = *addr
	return nil
}

func readKey(path *string, save *[]byte) error {

	if path != nil && len(*path) > 0 {

		key, err := ioutil.ReadFile(*path)
		if err != nil {
			return err
		}

		*save = key
	}

	return nil
}

// isValidAddress Проверка валидности адреса сервера
func isValidAddress(addr string) error {

	if addr == "" {
		return fmt.Errorf("address can not be empty")
	}

	parsedAddr := strings.Split(addr, ":")
	if len(parsedAddr) != 2 {
		return fmt.Errorf("need address in format host:port")
	}

	if len(parsedAddr[0]) > 0 && parsedAddr[0] != "localhost" {
		if ip := net.ParseIP(parsedAddr[0]); ip == nil {
			return fmt.Errorf("incorrect host: " + parsedAddr[0])
		}
	}

	if _, err := strconv.Atoi(parsedAddr[1]); err != nil {
		return fmt.Errorf("incorrect port: " + parsedAddr[1])
	}

	return nil
}
