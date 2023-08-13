package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand" // Do not use for crypt.
	"strings"
	"time"
)

const (
	TYPE_A   = 1
	CLASS_IN = 1
	// It is necessary to set any time for talking to a DNS resolver.
	// The encoding for the flags is defined in section 4.1.1 of RFC 1035.
	RECURSION_DESIRED = 1 << 8
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

func NewDNSHeader(id int, numQuestion int, flags int) DNSHeader {
	return DNSHeader{
		ID:             id,
		Flags:          flags,
		NumQuestion:    numQuestion,
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
			return nil, fmt.Errorf("encode error: %w", err)
		}

		err = binary.Write(&buf, binary.BigEndian, v)
		if err != nil {
			return nil, fmt.Errorf("encode error: %w", err)
		}
	}

	return buf.Bytes(), nil
}

func BuildQuery(domainName string, recordType int) (DNSHeader, DNSQuestion) {
	name, _ := EncodeDNSName(domainName)

	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))
	id := r.Intn(65535)

	header := NewDNSHeader(id, 1, RECURSION_DESIRED)
	question := DNSQuestion{
		Name:  name,
		Type:  recordType,
		Class: CLASS_IN,
	}

	return header, question
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
