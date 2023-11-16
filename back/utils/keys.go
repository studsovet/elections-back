package token

import (
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
