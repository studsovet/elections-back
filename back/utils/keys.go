package token

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

func ParsePrivateKey(privateKey string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privateKey))

	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func ParsePublicKey(publicKey string) (any, error) {
	block, _ := pem.Decode([]byte(publicKey))

	return x509.ParsePKIXPublicKey(block.Bytes)
}

func IsKeyMatched(publicKey any, privateKey *rsa.PrivateKey) bool {
	return privateKey.PublicKey.Equal(publicKey)
}

func DecryptWithPrivateKey(ciphertext []byte, priv *rsa.PrivateKey) ([]byte, error) {
	plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}
