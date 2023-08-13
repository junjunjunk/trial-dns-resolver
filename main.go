package main

import (
	"bytes"
)

type DNSHeader struct {
	ID             int
	Flags          int
	NumQuestion    int
	NumAnswers     int
	NumAuthorities int
	NumAdditionals int
}
type DNSQuestion struct {
	Name  []byte
	Type  int
	Class int
}

func NewDNSHeader(id int, flags int) DNSHeader {
	return DNSHeader{
		ID:             id,
		Flags:          flags,
		NumQuestion:    0,
		NumAnswers:     0,
		NumAuthorities: 0,
		NumAdditionals: 0,
	}
}

func main() {
	var buf bytes.Buffer

	// Network Packets usually uses BigEndian.
	// On the other hand, other situation uses LittleEndian.

}
