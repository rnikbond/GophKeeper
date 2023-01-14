package secret

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
)

func Encrypt(publicKey *rsa.PublicKey, data []byte) ([]byte, error) {
	if publicKey == nil {
		return data, nil
	}

	hashFunc := sha256.New()
	dataLen := len(data)
	step := publicKey.Size() - hashFunc.Size()*2 - 2
	var encryptedBytes []byte
	for start := 0; start < dataLen; start += step {
		finish := start + step
		if finish > dataLen {
			finish = dataLen
		}

		encryptedBlockBytes, err := rsa.EncryptOAEP(
			hashFunc,
			rand.Reader,
			publicKey,
			data[start:finish],
			nil)

		if err != nil {
			return nil, err
		}

		encryptedBytes = append(encryptedBytes, encryptedBlockBytes...)
	}

	return encryptedBytes, nil
}

func Decrypt(privKey *rsa.PrivateKey, data []byte) ([]byte, error) {

	if privKey == nil {
		return data, nil
	}

	dataLen := len(data)
	step := privKey.PublicKey.Size()
	var decryptedBytes []byte

	for start := 0; start < dataLen; start += step {
		finish := start + step
		if finish > dataLen {
			finish = dataLen
		}

		decryptedBlockBytes, err := privKey.Decrypt(nil, data[start:finish], &rsa.OAEPOptions{Hash: crypto.SHA256})
		if err != nil {
			return nil, err
		}

		decryptedBytes = append(decryptedBytes, decryptedBlockBytes...)
	}

	return decryptedBytes, nil
}

func GeneratePasswordHash(password, salt string) string {

	if len(password) < 1 {
		return password + salt
	}

	hash := sha256.New()
	hash.Write([]byte(password + salt))
	return hex.EncodeToString(hash.Sum(nil))
}
