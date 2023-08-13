package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
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

func EncodeDNSName(domainName string) ([]byte, error) {
	var buf bytes.Buffer
	tokenList := strings.Split(domainName, ".")

	for _, v := range tokenList {
		err := binary.Write(&buf, binary.BigEndian, len(v))
		if err != nil {
			fmt.Errorf("encode error: %w", err)
			return nil, err
		}

		err = binary.Write(&buf, binary.BigEndian, v)
		if err != nil {
			fmt.Errorf("encode error: %w", err)
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

func main() {
	var buf bytes.Buffer

	// Network Packets usually uses BigEndian.
	// On the other hand, other situation uses LittleEndian.
	err := binary.Write(&buf, binary.BigEndian)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
}
