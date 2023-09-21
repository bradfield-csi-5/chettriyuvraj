package main

type DNSHeader struct {
	ID           uint16
	Flags        uint16
	QCount       uint16
	AnsCount     uint16
	RRCount      uint16
	AddnlRRCount uint16
}

type DNSMessage struct {
	Header    DNSHeader
	Questions []DNSQuestion
	Answers   []DNSAnswer
}

type DNSQuestion struct {
	Name    []byte
	Type    uint16
	Class   uint16
	Namestr string /* Should have used a string - extended */
}

type DNSAnswer struct {
	Name     []byte
	Type     uint16
	Class    uint16
	TTL      uint32
	RDLength uint16
	RData    []byte
	Namestr  string /* Should have used a string - extended */
}
