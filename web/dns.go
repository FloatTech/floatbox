package web

import (
	"github.com/fumiama/terasu/dns"
)

func init() {
	dns.IPv4Servers.Add(&dns.DNSConfig{
		Servers: map[string][]string{
			"dot.360.cn": {
				"101.198.192.33:853",
				"112.65.69.15:853",
				"101.226.4.6:853",
				"218.30.118.6:853",
				"123.125.81.6:853",
				"140.207.198.6:853",
			},
		},
	})
	dns.IPv6Servers.Add(&dns.DNSConfig{
		Servers: map[string][]string{
			"dot.360.cn": {
				"101.198.192.33:853",
				"112.65.69.15:853",
				"101.226.4.6:853",
				"218.30.118.6:853",
				"123.125.81.6:853",
				"140.207.198.6:853",
			},
		},
	})
}
