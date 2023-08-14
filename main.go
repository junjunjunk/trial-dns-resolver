package main

import (
	"bytes"
	"fmt" // Do not use for crypt.

	"github.com/junjunjunk/trial-dns-resolver/client"
	"github.com/junjunjunk/trial-dns-resolver/model/dns"
	"github.com/junjunjunk/trial-dns-resolver/parser"
)

func main() {
	query, err := client.BuildQuery("www.example.com", dns.TYPE_A)
	if err != nil {
		fmt.Println(err)
		return
	}

	response, err := client.RequestDNSResolver(query)
	if err != nil {
		fmt.Println(err)
		return
	}

	reader := bytes.NewReader(response)

	packet, err := parser.ParseDNSPacket(reader)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf(packet.String())

}
