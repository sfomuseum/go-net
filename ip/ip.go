// Package net provides network-related methods.
package ip

import (
	"encoding/binary"
	"fmt"
	_ "log"
	"net"
	"net/http"
	"strings"
)

// Derive a numeric IP address from 'req'.
func DeriveAddress(req *http.Request) (uint32, error) {

	remote_addr := req.RemoteAddr

	// because net.ParseIP can't parse stuff like "127.0.0.1:5656373"
	// which is what remote_addr will be if we're running in debug mode
	// on localhost (20200817/thisisaaronland)

	if strings.HasPrefix(remote_addr, "127.0.0.1") {
		remote_addr = "127.0.0.1"
	}

	ip_addr := net.ParseIP(remote_addr)

	if ip_addr == nil {
		return 0, fmt.Errorf("Failed to parse IP address")
	}

	return IPToInt(ip_addr), nil
}

// IPToInt converts ip_addr in to a `uint32` value.
func IPToInt(ip_addr net.IP) uint32 {

	/*
		if len(ip_addr) == 16 {
			return binary.BigEndian.Uint32(ip_addr[12:16])
		}
	*/

	return binary.BigEndian.Uint32(ip_addr.To4())
}

// IntToIP converts 'n' in to a `net.IP` instance.
func IntToIP(n uint32) net.IP {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, n)
	return ip
}
