package ip

import (
	"net"
	"testing"
)

func TestIPToInt(t *testing.T) {

	tests := map[string]uint32{
		"127.0.0.1": 2130706433,
		"8.8.8.8":   134744072,
	}

	for str_ip, expected := range tests {

		ip := net.ParseIP(str_ip)
		i := IPToInt(ip)

		if i != expected {
			t.Fatalf("Unexpected value for '%s': %d", str_ip, i)
		}
	}
}

func TestIntToIP(t *testing.T) {

	tests := map[uint32]string{
		2130706433: "127.0.0.1",
		134744072:  "8.8.8.8",
	}

	for i, expected := range tests {

		ip := IntToIP(i)

		if ip.String() != expected {
			t.Fatalf("Unexpected value for '%d': %s", i, ip.String())
		}
	}
}
