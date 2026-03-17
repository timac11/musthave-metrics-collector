package middleware

import (
	"github.com/timac11/musthave-metrics-collector/internal/encryption"
)

type Middleware struct {
	signingKey string
	decoder    *encryption.Decoder
}

func NewMiddleware(hashingKey, pkPath string) (*Middleware, error) {
	if pkPath != "" {
		decoder, err := encryption.NewDecoder(pkPath)

		if err != nil {
			return nil, err
		}

		return &Middleware{signingKey: hashingKey, decoder: decoder}, nil
	}

	return &Middleware{signingKey: hashingKey}, nil
}
