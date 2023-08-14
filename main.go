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

	buf := bytes.NewReader(response)

	respHeader, err := parser.ParseHeader(buf)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%+v\n", respHeader)

	respQuestion, err := parser.ParseQuestion(buf)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%+v\n", respQuestion)

	respRecord, err := parser.ParseRecord(buf)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%+v\n", respRecord)

}
