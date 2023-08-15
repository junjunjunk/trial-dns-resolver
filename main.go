package main

import (
	"fmt"

	"github.com/junjunjunk/trial-dns-resolver/client"
	"github.com/junjunjunk/trial-dns-resolver/model/dns"
)

func main() {
	packet, err := client.SendQuery("198.41.0.4", "google.com", dns.TYPE_TXT)
	if err != nil {
		return
	}
	fmt.Println(packet.String())
}
