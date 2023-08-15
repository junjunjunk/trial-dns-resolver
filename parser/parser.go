package parser

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/junjunjunk/trial-dns-resolver/model/dns"
)

func ParseHeader(reader *bytes.Reader) (*dns.DNSHeader, error) {
	var header dns.DNSHeader
	//  Fields is a 2-byte integer, so there are 12 bytes in all to read.
	err := binary.Read(reader, binary.BigEndian, &header)
	if err != nil {
		return nil, err
	}

	return &header, err
}

func ParseQuestion(reader *bytes.Reader) (*dns.DNSQuestion, error) {
	var question dns.DNSQuestion
	//  Fields is a 2-byte integer, so there are 12 bytes in all to read.
	question.Name = DecodeName(reader)
	err := binary.Read(reader, binary.BigEndian, &question.Type)
	if err != nil {
		return nil, err
	}
	err = binary.Read(reader, binary.BigEndian, &question.Class)
	if err != nil {
		return nil, err
	}
	return &question, err
}

func ParseRecord(reader *bytes.Reader) (*dns.DNSRecord, error) {

	var record dns.DNSRecord

	record.Name = DecodeName(reader)

	err := binary.Read(reader, binary.BigEndian, &record.Type)
	if err != nil {
		return nil, err
	}

	err = binary.Read(reader, binary.BigEndian, &record.Class)
	if err != nil {
		return nil, err
	}

	err = binary.Read(reader, binary.BigEndian, &record.TTL)
	if err != nil {
		return nil, err
	}

	_, err = reader.Read(make([]byte, 2)) // skip dataLen
	if err != nil {
		return nil, err
	}

	if record.Type == dns.TYPE_NS {
		record.Data = DecodeName(reader)
	} else if record.Type == dns.TYPE_A {
		err = binary.Read(reader, binary.BigEndian, &record.Data)
		if err != nil {
			return nil, err
		}
	} else {
		err = binary.Read(reader, binary.BigEndian, &record.Data)
		if err != nil {
			return nil, err
		}
	}

	return &record, nil
}

// https://datatracker.ietf.org/doc/html/rfc1035#section-4.1.4
func DecodeName(reader *bytes.Reader) []byte {

	var parts [][]byte
	for {
		lengthByte, err := reader.ReadByte()
		// if err or 0x00
		if (err != nil) || (lengthByte == 0) {
			break
		}

		length := int(lengthByte)
		// The maximum length of each part is 63! The first 2 bits of the byte 192 (11000000 in binary) are 11,
		// and any length that starts with the bits 11 is code for “this is compressed”.
		// https://datatracker.ietf.org/doc/html/rfc1035#section-4.1.4
		if length&0b1100_0000 != 0 {
			compressedPart := decodeCompressedName(length, reader)
			parts = append(parts, compressedPart)
			break
		} else {
			part := make([]byte, length)
			_, err = reader.Read(part)
			if err != nil {
				break
			}
			parts = append(parts, part)
		}

	}
	return bytes.Join(parts, []byte("."))
}

func decodeCompressedName(length int, reader *bytes.Reader) []byte {
	pointerByte, err := reader.ReadByte()
	if err != nil {
		return nil
	}
	pointer := byte(length&0b0011_1111) + pointerByte
	currentPosition, err := reader.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil
	}

	reader.Seek(int64(pointer), io.SeekStart)
	result := DecodeName(reader)
	reader.Seek(currentPosition, io.SeekStart)
	return result
}

func ParseDNSPacket(reader *bytes.Reader) (*dns.DNSPacket, error) {
	dnsPacket := &dns.DNSPacket{}
	header, err := ParseHeader(reader)
	if err != nil {
		return nil, err
	}
	dnsPacket.Header = header
	for i := 0; i < int(dnsPacket.Header.NumQuestions); i++ {
		question, err := ParseQuestion(reader)
		if err != nil {
			return nil, err
		}
		dnsPacket.Questions = append(dnsPacket.Questions, question)
	}

	for i := 0; i < int(dnsPacket.Header.NumAnswers); i++ {
		answer, err := ParseRecord(reader)
		if err != nil {
			return nil, err
		}
		dnsPacket.Answers = append(dnsPacket.Answers, answer)
	}

	for i := 0; i < int(dnsPacket.Header.NumAuthorities); i++ {
		auhtority, err := ParseRecord(reader)
		if err != nil {
			return nil, err
		}
		dnsPacket.Authorities = append(dnsPacket.Authorities, auhtority)
	}

	for i := 0; i < int(dnsPacket.Header.NumAdditionals); i++ {
		additinonal, err := ParseRecord(reader)
		if err != nil {
			return nil, err
		}
		dnsPacket.Additionals = append(dnsPacket.Additionals, additinonal)
	}

	return dnsPacket, nil
}
