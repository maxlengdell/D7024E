package main

import (
	"fmt"
	"net"

	d7024e "github.com/maxlengdell/D7024E/d7024e"
)

func main() {
	fmt.Println("hello")
	d7024e.Listen("0.0.0.0", 8080)
	//name, err := os.Hostname()
	ip := getadress()

}
func getadress() string {
	iface, _ := net.InterfaceByName("eth0")
	addr, _ := iface.Addrs()
	for _, addr := range addr {
		var ip net.IP
		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		}
		return ip.String()
		// process IP address
	}
	return "Could not get local ip address"
}
