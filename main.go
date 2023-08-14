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
	NumQuestions   int
	NumAnswers     int
	NumAuthorities int
	NumAdditionals int
}
type DNSQuestion struct {
	Name  []byte
	Type  int
	Class int
}

func NewDNSHeader(id int, numQuestions int, flags int) DNSHeader {
	return DNSHeader{
		ID:             id,
		Flags:          flags,
		NumQuestions:   numQuestions,
		NumAnswers:     0,
		NumAuthorities: 0,
		NumAdditionals: 0,
	}
}

func EncodeDNSName(domainName string) ([]byte, error) {
	var buf bytes.Buffer
	tokenList := strings.Split(domainName, ".")

	for _, v := range tokenList {
		err := binary.Write(&buf, binary.BigEndian, byte(len(v)))
		if err != nil {
			return nil, fmt.Errorf("encode error: %w", err)
		}

		err = binary.Write(&buf, binary.BigEndian, []byte(v))
		if err != nil {
			return nil, fmt.Errorf("encode error: %w", err)
		}
	}
	// TODO?: Add End of Bytes 0x00.
	return buf.Bytes(), nil
}

func HeaderToBytes(header DNSHeader) ([]byte, error) {
	var buf bytes.Buffer

	err := binary.Write(&buf, binary.BigEndian, header.ID)
	if err != nil {
		return nil, fmt.Errorf("headerToBytes error: %w", err)
	}

	err = binary.Write(&buf, binary.BigEndian, header.Flags)
	if err != nil {
		return nil, fmt.Errorf("headerToBytes error: %w", err)
	}

	err = binary.Write(&buf, binary.BigEndian, byte(header.NumQuestions))
	if err != nil {
		return nil, fmt.Errorf("headerToBytes error: %w", err)
	}

	err = binary.Write(&buf, binary.BigEndian, byte(header.NumAnswers))
	if err != nil {
		return nil, fmt.Errorf("headerToBytes error: %w", err)
	}

	err = binary.Write(&buf, binary.BigEndian, byte(header.NumAuthorities))
	if err != nil {
		return nil, fmt.Errorf("headerToBytes error: %w", err)
	}

	err = binary.Write(&buf, binary.BigEndian, byte(header.NumAdditionals))
	if err != nil {
		return nil, fmt.Errorf("headerToBytes error: %w", err)
	}

	return buf.Bytes(), nil
}

func QuestionToBytes(header DNSHeader) ([]byte, error) {
	var buf bytes.Buffer
	return buf.Bytes(), nil
}

func BuildQuery(domainName string, recordType int) (DNSHeader, DNSQuestion) {
	name, err := EncodeDNSName(domainName)
	if err != nil {
		fmt.Println(err)
	}

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

	header, question := BuildQuery("www.example.com", TYPE_A)
	fmt.Println(header)
	fmt.Println(question)
	// Network Packets usually uses BigEndian.
	// On the other hand, other situation uses LittleEndian.
	// TODO: err handling
	binary.Write(&buf, binary.BigEndian, header)
	binary.Write(&buf, binary.BigEndian, question)

	fmt.Printf("%q\n", buf)
}
