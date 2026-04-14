package util

import (
	"net"
)

// GetOutboundIP get preferred outbound ip of the client
func GetOutboundIP(address string) (string, error) {
	conn, err := net.Dial("udp", address)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String(), nil
}
