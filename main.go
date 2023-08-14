package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/rand" // Do not use for crypt.
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

type DNSHeader struct {
	ID             uint16
	Flags          uint16
	NumQuestions   uint16
	NumAnswers     uint16
	NumAuthorities uint16
	NumAdditionals uint16
}

// REF: https://datatracker.ietf.org/doc/html/rfc1035#section-4.1.2
type DNSQuestion struct {
	Name  []byte
	Type  uint16
	Class uint16
}

func NewDNSHeader(id uint16, numQuestions uint16, flags uint16) DNSHeader {
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

	err := buf.WriteByte(0x00)
	if err != nil {
		fmt.Println("Error writing null terminator:", err)
		return nil, err
	}
	return buf.Bytes(), nil
}

// Fixed Length Encoding
func HeaderToBytes(header DNSHeader) ([]byte, error) {
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
func QuestionToBytes(question DNSQuestion) ([]byte, error) {
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

	header := NewDNSHeader(id, 1, RECURSION_DESIRED)
	question := DNSQuestion{
		Name:  name,
		Type:  recordType,
		Class: CLASS_IN,
	}

	headerBytes, err := HeaderToBytes(header)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	questionBytes, err := QuestionToBytes(question)
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

func main() {
	query, err := BuildQuery("www.example.com", TYPE_A)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("query:", hex.EncodeToString(query))

	// port53: dns port
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	_, err = conn.Write(query)
	if err != nil {
		panic(err)
	}

	response := make([]byte, 1024)
	_, err = conn.Read(response)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%q\n", response)
}
