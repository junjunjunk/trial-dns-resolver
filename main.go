package main

import (
	"fmt"

	"github.com/junjunjunk/trial-dns-resolver/client"
)

func main() {
	ip, err := client.LookUpDomain("www.example.com")
	if err != nil {
		return
	}
	fmt.Println(ip)
}
