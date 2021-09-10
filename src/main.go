package main

import (
	"net"

	d7024e "github.com/maxlengdell/D7024E/d7024e"
)

func main() {
	localIP := getadress()
	firstIP := "172.18.0.2"
	var kademlia *d7024e.Kademlia
	if localIP == firstIP {
		//If first node in the network, bootstrap
		//Enter bootstrap sequence
		kademlia = d7024e.Bootstrap(localIP, 8080)
		//rt := d7024e.NewRoutingTable(*myContact)

	} else {
		kademlia = d7024e.JoinNetwork(firstIP, localIP, 8080)
		//rt := d7024e.NewRoutingTable(*myContact)

	}
	msgChan := make(chan d7024e.InternelMessage)
	go kademlia.Net.Listen(localIP, 8080, msgChan)
	go kademlia.Net.Listen(localIP, 1010, msgChan)

	kademlia.HandleMessage(msgChan)

	//id := d7024e.NewRandomKademliaID()
	//contact := d7024e.NewContact(id, "172.18.0.2")
	//net.SendPingMessage(&contact)

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
