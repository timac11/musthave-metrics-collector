package middleware

import (
	"net/netip"

	"github.com/timac11/musthave-metrics-collector/internal/common/encryption"
)

type Middleware struct {
	signingKey string
	decoder    *encryption.Decoder
	subnet     *netip.Prefix
}

func NewMiddleware(hashingKey, pkPath, trustedSubnet string) (*Middleware, error) {
	var decoder *encryption.Decoder
	var subnet *netip.Prefix
	var err error

	if pkPath != "" {
		decoder, err = encryption.NewDecoder(pkPath)

		if err != nil {
			return nil, err
		}
	}

	if trustedSubnet != "" {
		parsed, err := netip.ParsePrefix(trustedSubnet)

		if err != nil {
			return nil, err
		}

		subnet = &parsed
	}

	return &Middleware{signingKey: hashingKey, decoder: decoder, subnet: subnet}, nil
}
