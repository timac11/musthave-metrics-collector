package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

// CalculateSignature calculate signature of serializable structure
func CalculateSignature(data any, signingKey string) (string, error) {
	if signingKey == "" {
		return "", nil
	}

	var bytesToSign []byte

	if dataBytes, ok := data.([]byte); ok {
		bytesToSign = dataBytes
	} else {
		dataBytes, err := json.Marshal(data)
		if err != nil {
			return "", fmt.Errorf("failed parse object to sign: %w", err)
		}
		bytesToSign = dataBytes
	}

	h := hmac.New(sha256.New, []byte(signingKey))
	_, err := h.Write(bytesToSign)

	if err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
