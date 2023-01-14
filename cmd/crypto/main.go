package main

import (
	"GophKeeper/pkg/logzap"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"go.uber.org/zap"
	"io/ioutil"
)

func GenerateRsaKeyPair() (*rsa.PrivateKey, *rsa.PublicKey) {
	key, _ := rsa.GenerateKey(rand.Reader, 4096)
	return key, &key.PublicKey
}

func PrivateToString(key *rsa.PrivateKey) string {
	bytes := x509.MarshalPKCS1PrivateKey(key)
	private := pem.EncodeToMemory(
		&pem.Block{

			Type:  "RSA PRIVATE KEY",
			Bytes: bytes,
		},
	)
	return string(private)
}

func PublicToString(key *rsa.PublicKey) (string, error) {
	bytes, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return ``, err
	}

	public := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: bytes,
		},
	)

	return string(public), nil
}

func ExportToFile(data string, file string) error {
	return ioutil.WriteFile(file, []byte(data), 0644)
}

func main() {

	logzap.ConfigZapLogger()
	logger := zap.L()

	privatePath := "private.key"
	publicPath := "public.key"

	privateKey, publicKey := GenerateRsaKeyPair()

	privatePEM := PrivateToString(privateKey)
	publicPEM, _ := PublicToString(publicKey)

	if err := ExportToFile(privatePEM, privatePath); err != nil {
		logger.Error("failed export to file private key", zap.Error(err))
	} else {
		logger.Info("success export private key to file", zap.String("file", privatePath))
	}

	if err := ExportToFile(publicPEM, publicPath); err != nil {
		logger.Error("failed export to file public key", zap.Error(err))
	} else {
		logger.Info("success export public key to file", zap.String("file", publicPath))
	}
}
