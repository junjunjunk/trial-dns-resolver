package dns

import "fmt"

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
	Data  [4]byte
}

type DNSPacket struct {
	Header      *DNSHeader
	Questions   []*DNSQuestion
	Answers     []*DNSRecord
	Authorities []*DNSRecord
	Additionals []*DNSRecord
}

func (p *DNSPacket) String() string {
	return fmt.Sprintf("Header: %+v\nQuestions: %s\nAnswers: %s\nAdditionals: %s\n",
		p.Header, formatQuestions(p.Questions), formatRecords(p.Answers), formatRecords(p.Additionals))
}

func formatQuestions(questions []*DNSQuestion) string {
	var result string
	if len(questions) == 0 {
		return "x"
	}
	for i, q := range questions {
		if i > 0 {
			fmt.Print(", ")
		}
		result += fmt.Sprintf("%+v", (q))
	}
	return result
}

func formatRecords(records []*DNSRecord) string {
	var result string
	if len(records) == 0 {
		return "x"
	}
	for i, r := range records {
		if i > 0 {
			fmt.Print(", ")
		}
		result += fmt.Sprintf("%+v", (r))
		fmt.Printf("%q\n", r)

	}
	return result
}
