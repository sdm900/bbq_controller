package network

import (
	"net"
	"os"
)

func GetLocalIP() string {
	// returns the first sensible IP address
	addrs, e := net.InterfaceAddrs()
	if e != nil {
		return "Unknown"
	}

	var ip string
	for _, addr := range addrs {
		ipnet, ok := addr.(*net.IPNet)
		if ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil && !ipnet.IP.IsMulticast() {
			ip = ipnet.IP.To4().String()

			break
		}
	}

	return ip
}

func Hostname() string {
	hn, e := os.Hostname()
	if e != nil || hn == "" {
		hn = "localhost"
	}

	return hn
}
