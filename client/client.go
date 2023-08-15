package client

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand" // Do not use for crypt.
	"net"
	"strings"
	"time"

	"github.com/junjunjunk/trial-dns-resolver/model/dns"
	"github.com/junjunjunk/trial-dns-resolver/parser"
)

func newDNSHeader(id uint16, numQuestions uint16, flags uint16) dns.DNSHeader {
	return dns.DNSHeader{
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
		return nil, fmt.Errorf("Error writing null terminator: %w", err)
	}
	return buf.Bytes(), nil
}

// Fixed Length Encoding
func headerToBytes(header dns.DNSHeader) ([]byte, error) {
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
func questionToBytes(question dns.DNSQuestion) ([]byte, error) {
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

func buildQuery(domainName string, recordType uint16) ([]byte, error) {
	name, err := EncodeDNSName(domainName)

	if err != nil {
		return nil, fmt.Errorf("buildQuery error: %w", err)
	}

	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))
	id := uint16(r.Intn(65535))

	// header := newDNSHeader(id, 1, dns.RECURSION_DESIRED)
	header := newDNSHeader(id, 1, 0)
	question := dns.DNSQuestion{
		Name:  name,
		Type:  recordType,
		Class: dns.CLASS_IN,
	}

	headerBytes, err := headerToBytes(header)
	if err != nil {
		return nil, fmt.Errorf("buildQuery error: %w", err)
	}
	questionBytes, err := questionToBytes(question)
	if err != nil {
		return nil, fmt.Errorf("buildQuery error: %w", err)
	}

	var buf bytes.Buffer
	_, err = buf.Write(headerBytes)
	if err != nil {
		return nil, fmt.Errorf("buf write error: %w", err)
	}

	_, err = buf.Write(questionBytes)
	if err != nil {
		return nil, fmt.Errorf("buf write error: %w", err)
	}

	return buf.Bytes(), nil
}

func requestDNSResolver(ipAddress string, query []byte) ([]byte, error) {
	// port53: dns port
	conn, err := net.Dial("udp", ipAddress+":53")
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

func LookUpDomain(domainName string) (string, error) {
	query, err := buildQuery(domainName, dns.TYPE_A)
	if err != nil {
		return "", fmt.Errorf("build query error: %w", err)
	}

	response, err := requestDNSResolver("8.8.8.8", query)
	if err != nil {
		return "", fmt.Errorf("request dns resolver error: %w", err)
	}

	reader := bytes.NewReader(response)

	packet, err := parser.ParseDNSPacket(reader)
	if err != nil {
		return "", fmt.Errorf("parse packet error: %w", err)
	}

	return packet.IP(), nil
}

func SendQuery(ipAddress string, domainName string, recordType uint16) (*dns.DNSPacket, error) {
	query, err := buildQuery(domainName, recordType)
	if err != nil {
		return nil, fmt.Errorf("build query error: %w", err)
	}

	response, err := requestDNSResolver(ipAddress, query)
	if err != nil {
		return nil, fmt.Errorf("request dns resolver error: %w", err)
	}

	reader := bytes.NewReader(response)

	packet, err := parser.ParseDNSPacket(reader)
	if err != nil {
		return nil, fmt.Errorf("parse packet error: %w", err)
	}

	return packet, nil
}

func Resolve(domainName string, recordType uint16) (string, error) {
	nameServer := "198.41.0.4"
	for {
		fmt.Printf("Querying %s for %s\n", nameServer, domainName)
		packet, err := SendQuery(nameServer, domainName, recordType)
		if err != nil {
			return "", err
		}

		if ip := packet.GetAnswer(); ip != "" {
			return ip, nil
		} else if nsIP := packet.GetNameserverIP(); nsIP != "" {
			nameServer = nsIP
		} else if nsDomain := packet.GetNameserver(); nsDomain != "" {
			nameServer, err = Resolve(nsDomain, dns.TYPE_A)
			if err != nil {
				return "", err
			}
		} else {
			return "", fmt.Errorf("something went wrong")
		}
	}
}
