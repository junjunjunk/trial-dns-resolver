package dns

import (
	"fmt"
)

const (
	TYPE_A   = 1
	TYPE_NS  = 2
	TYPE_TXT = 16
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

// name: the domain name
// type_: A, AAAA, MX, NS, TXT, etc (encoded as an integer)
// class: always the same (1). We’ll ignore this.
// ttl: how long to cache the query for. We’ll ignore this.
// data: the record’s content, like the IP address.
type DNSRecord struct {
	Name  []byte
	Type  uint16
	Class uint16
	TTL   uint32
	Data  []byte
}

type DNSPacket struct {
	Header      *DNSHeader
	Questions   []*DNSQuestion
	Answers     []*DNSRecord
	Authorities []*DNSRecord
	Additionals []*DNSRecord
}

func (p *DNSPacket) String() string {
	return fmt.Sprintf("[Header]\n%+v\n[Questions]\n%s\n[Answers]\n%s\n[Authorities]\n%s\n[Additionals]\n%s\n",
		p.Header, formatQuestions(p.Questions), formatRecords(p.Answers), formatRecords(p.Authorities), formatRecords(p.Additionals))
}

func (p *DNSPacket) IP() string {
	var ipAddresses [][]byte
	for _, a := range p.Answers {
		ipAddresses = append(ipAddresses, a.Data[:])
	}

	if len(ipAddresses) == 0 {
		return ""
	}

	var result string
	result = "["
	for i, ip := range ipAddresses {
		if i > 0 {
			result += fmt.Sprintf(", ")
		}
		for j, b := range ip {
			if j > 0 {
				result += fmt.Sprintf(".")
			}
			result += fmt.Sprintf("%d", b)
		}
	}
	result += "]"

	return result
}

func formatQuestions(questions []*DNSQuestion) string {
	var result string
	if len(questions) == 0 {
		return "x"
	}
	for _, q := range questions {

		result += fmt.Sprintf("%+v", (q))
	}
	return result
}

func formatRecords(records []*DNSRecord) string {
	var result string
	if len(records) == 0 {
		return "x"
	}
	for _, r := range records {
		result += fmt.Sprintf("%+v\n", (r))
	}
	return result
}

func (p *DNSPacket) GetNameserver() string {
	for _, x := range p.Authorities {
		if x.Type == TYPE_NS {
			return string(x.Data)
		}
	}
	return ""
}

func (p *DNSPacket) GetNameserverIP() string {
	for _, x := range p.Additionals {
		if x.Type == TYPE_A {
			result := ""
			for i, b := range x.Data {
				if i > 0 {
					result += fmt.Sprintf(".")
				}
				result += fmt.Sprintf("%d", b)
			}
			return result
		}
	}
	return ""
}

func (p *DNSPacket) GetAnswer() string {
	for _, x := range p.Answers {
		if x.Type == TYPE_A {
			result := ""
			for i, b := range x.Data {
				if i > 0 {
					result += fmt.Sprintf(".")
				}
				result += fmt.Sprintf("%d", b)
			}
			return result
		}
	}
	return ""
}
