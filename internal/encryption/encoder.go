package encryption

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
)

type Encoder struct {
	rec *x509.Certificate
}

func NewEncoder(path string) (*Encoder, error) {
	certificateBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	certificatePemBlock, _ := pem.Decode(certificateBytes)
	if certificatePemBlock == nil {
		return nil, errors.New("certificate not found")
	}

	certificate, err := x509.ParseCertificate(certificatePemBlock.Bytes)
	if err != nil {
		return nil, err
	}

	return &Encoder{rec: certificate}, nil
}

func (e *Encoder) Encode(message []byte) ([]byte, error) {
	encryptedMessage, err := rsa.EncryptPKCS1v15(rand.Reader, e.rec.PublicKey.(*rsa.PublicKey), message)
	if err != nil {
		return []byte{}, err
	}

	return encryptedMessage, nil
}
