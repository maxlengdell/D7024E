package main

import (
	"net"

	"github.com/maxlengdell/D7024E/d7024e"
)

func main() {
	net := d7024e.Network{}
	localIP := getadress()
	id := d7024e.NewRandomKademliaID()

	go d7024e.Listen(localIP, 8080)
	contact := d7024e.NewContact(id, "172.18.0.2")
	net.SendPingMessage(&contact)

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
