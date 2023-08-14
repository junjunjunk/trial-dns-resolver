package parser

import (
	"bytes"
	"encoding/binary"

	"github.com/junjunjunk/trial-dns-resolver/model/dns"
)

func ParseHeader(reader *bytes.Buffer) (*dns.DNSHeader, error) {
	var header dns.DNSHeader
	//  Fields is a 2-byte integer, so there are 12 bytes in all to read.
	err := binary.Read(reader, binary.BigEndian, &header)
	if err != nil {
		return nil, err
	}

	return &header, err
}
