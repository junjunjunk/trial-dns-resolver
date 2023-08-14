package client

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
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

type dnsHeader struct {
	ID             uint16
	Flags          uint16
	NumQuestions   uint16
	NumAnswers     uint16
	NumAuthorities uint16
	NumAdditionals uint16
}

// REF: https://datatracker.ietf.org/doc/html/rfc1035#section-4.1.2
type dnsQuestion struct {
	Name  []byte
	Type  uint16
	Class uint16
}

func newDNSHeader(id uint16, numQuestions uint16, flags uint16) dnsHeader {
	return dnsHeader{
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

	err := buf.WriteByte(0x00)
	if err != nil {
		fmt.Println("Error writing null terminator:", err)
		return nil, err
	}
	return buf.Bytes(), nil
}

// Fixed Length Encoding
func headerToBytes(header dnsHeader) ([]byte, error) {
	var buf bytes.Buffer

	// decodedData, _ := hex.DecodeString("0x8298")
	err := binary.Write(&buf, binary.BigEndian, header.ID)
	if err != nil {
		return nil, fmt.Errorf("headerToBytes error: %w", err)
	}

	err = binary.Write(&buf, binary.BigEndian, header.Flags)
	if err != nil {
		return nil, fmt.Errorf("headerToBytes error: %w", err)
	}

	err = binary.Write(&buf, binary.BigEndian, header.NumQuestions)
	if err != nil {
		return nil, fmt.Errorf("headerToBytes error: %w", err)
	}

	err = binary.Write(&buf, binary.BigEndian, header.NumAnswers)
	if err != nil {
		return nil, fmt.Errorf("headerToBytes error: %w", err)
	}

	err = binary.Write(&buf, binary.BigEndian, header.NumAuthorities)
	if err != nil {
		return nil, fmt.Errorf("headerToBytes error: %w", err)
	}

	err = binary.Write(&buf, binary.BigEndian, header.NumAdditionals)
	if err != nil {
		return nil, fmt.Errorf("headerToBytes error: %w", err)
	}

	return buf.Bytes(), nil
}

// Fixed Length Encoding
func questionToBytes(question dnsQuestion) ([]byte, error) {
	var buf bytes.Buffer

	err := binary.Write(&buf, binary.BigEndian, question.Name)
	if err != nil {
		return nil, fmt.Errorf("headerToBytes error: %w", err)
	}

	err = binary.Write(&buf, binary.BigEndian, question.Type)
	if err != nil {
		return nil, fmt.Errorf("headerToBytes error: %w", err)
	}

	err = binary.Write(&buf, binary.BigEndian, question.Class)
	if err != nil {
		return nil, fmt.Errorf("headerToBytes error: %w", err)
	}

	return buf.Bytes(), nil
}

func BuildQuery(domainName string, recordType uint16) ([]byte, error) {
	name, err := EncodeDNSName(domainName)

	if err != nil {
		fmt.Println(err)
	}

	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))
	id := uint16(r.Intn(65535))

	header := newDNSHeader(id, 1, RECURSION_DESIRED)
	question := dnsQuestion{
		Name:  name,
		Type:  recordType,
		Class: CLASS_IN,
	}

	headerBytes, err := headerToBytes(header)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	questionBytes, err := questionToBytes(question)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var buf bytes.Buffer
	_, err = buf.Write(headerBytes)
	if err != nil {
		fmt.Println(fmt.Errorf("buf write error: %w", err))
		return nil, err
	}

	_, err = buf.Write(questionBytes)
	if err != nil {
		fmt.Println(fmt.Errorf("buf write error: %w", err))
		return nil, err
	}

	return buf.Bytes(), nil
}

func RequestDNSResolver(query []byte) ([]byte, error) {
	// port53: dns port
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		return nil, fmt.Errorf("net dial error: %w", err)
	}
	defer conn.Close()

	_, err = conn.Write(query)
	if err != nil {
		return nil, fmt.Errorf("conn write error: %w", err)
	}

	response := make([]byte, 1024)
	_, err = conn.Read(response)
	if err != nil {
		return nil, fmt.Errorf("conn read error: %w", err)
	}

	return response, nil
}
