package main

//Docker attach shell
//docker exec -it upbeat_dewdney /bin/sh
import (
	"fmt"
	"net"
	"strings"
	"time"

	d7024e "github.com/maxlengdell/D7024E/d7024e"
)

func main() {
	localIP := getadress()
	firstIP := localIP[:len(localIP)-2] + ".2"
	var kademlia *d7024e.Kademlia
	msgChan := make(chan d7024e.InternalMessage)
	go d7024e.Listen(localIP, 8080, msgChan) //External comm
	go d7024e.Listen(localIP, 1010, msgChan) //CLI

	if localIP == firstIP {
		fmt.Println("Bootstrap node")
		//Enter bootstrap sequence
		kademlia = d7024e.Bootstrap(localIP, 8080)

	} else {
		time.Sleep(5 * time.Second)
		kademlia = d7024e.JoinNetwork(firstIP, localIP, 8080)
	}
	kademlia.HandleMessage(msgChan)
}

func getadress() string {
	//iface, _ := net.InterfaceByName("eth0")
	ifaces, _ := net.Interfaces()
	for _, iface := range ifaces {
		if strings.HasPrefix(iface.Name, "lo") { // TODO: Figure out how to test iface.Flags
			continue // Skip loopback
		}
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
	}
	return "Could not get local ip address"
}
