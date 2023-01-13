package server

import (
	"flag"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type Config struct {
	AddrGRPC    string `env:"ADDRESS_RPC" json:"address_rpc"`
	SecretKey   string `env:"SECRET_KEY"  json:"secret_key"`
	DatabaseURI string `env:"DatabaseURI" json:"database_uri"`
}

// NewConfig Конфигурация сервера
func NewConfig() *Config {

	return &Config{
		AddrGRPC:    ":3200",
		DatabaseURI: "user=postgres password=postgres dbname=GophKeeper sslmode=disable",
	}
}

// ParseArgs Разброр аргументов командной строки
func (cfg *Config) ParseArgs() error {

	addr := flag.String("a", "", "address grpc gate")
	secret := flag.String("s", "", "secret key for JWT")
	dsn := flag.String("d", "", "database DSN")
	flag.Parse()

	if addr == nil || len(*addr) == 0 {
		*addr = cfg.AddrGRPC
	}

	if err := isValidAddress(*addr); err != nil {
		return err
	}

	cfg.AddrGRPC = *addr

	if dsn != nil && len(*dsn) > 0 {
		cfg.DatabaseURI = *dsn
	}

	if secret != nil {
		cfg.SecretKey = *secret
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
