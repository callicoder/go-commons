package requestutil

import (
	"net"
	"net/http"
	"strings"
)

func GetIpAddress(r *http.Request) net.IP {
	address := ""

	header := r.Header.Get("X-Forwarded-For")
	if len(header) > 0 {
		addresses := strings.Fields(header)
		if len(addresses) > 0 {
			address = strings.TrimRight(addresses[0], ",")
		}
	}

	if len(address) == 0 {
		address = r.Header.Get("X-Real-IP")
	}

	if len(address) == 0 {
		address, _, _ = net.SplitHostPort(r.RemoteAddr)
	}

	userIP := net.ParseIP(address)
	return userIP
}
