package main

import (
	"fmt"

	"github.com/junjunjunk/trial-dns-resolver/client"
	"github.com/junjunjunk/trial-dns-resolver/model/dns"
)

func main() {
	// 198.41.0.4
	// Real DNS resolvers actually do hardcode the IP addresses of the root nameservers. This is because if you’re implementing DNS, you have to start somewhere
	packet, err := client.SendQuery("192.12.94.30", "google.com", dns.TYPE_TXT)
	if err != nil {
		return
	}
	fmt.Println(packet.String())
}
