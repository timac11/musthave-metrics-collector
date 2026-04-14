package encryption

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
)

type Decoder struct {
	pk *rsa.PrivateKey
}

func NewDecoder(path string) (*Decoder, error) {
	privateKeyBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	privateKeyPemBlock, _ := pem.Decode(privateKeyBytes)
	if privateKeyPemBlock == nil {
		return nil, errors.New("private key not found")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyPemBlock.Bytes)
	if err != nil {
		return nil, err
	}

	return &Decoder{pk: privateKey}, nil
}

func (d *Decoder) Decode(message []byte) ([]byte, error) {
	decryptedMessage, err := rsa.DecryptPKCS1v15(rand.Reader, d.pk, message)
	if err != nil {
		return []byte{}, err
	}

	return decryptedMessage, nil
}
