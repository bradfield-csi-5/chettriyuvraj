package main

type DNSHeader struct {
	ID           uint16
	Flags        uint16
	QCount       uint16
	AnsCount     uint16
	RRCount      uint16
	AddnlRRCount uint16
}

type DNSQuery struct {
	Header    DNSHeader
	Questions []DNSQuestion
}

type DNSMessage struct {
	Header    DNSHeader
	Questions []DNSQuestion
	Answers   []DNSAnswer
}

type DNSQuestion struct {
	Name  []byte
	Type  uint16
	Class uint16
}

type DNSAnswer struct {
	Name     []byte
	Type     uint16
	Class    uint16
	TTL      uint32
	RDLength uint16
	RData    []byte
}

type IPv4Packet struct {
	Version        uint8
	IHL            uint8
	DSCP           uint8
	ECN            uint8
	TotalLen       uint16
	Id             uint16
	Flags          uint8
	FragmentOffset uint16
	TTL            uint8
	Protocol       uint8
	HeaderChecksum uint16
	SourceIP       uint32
	DestIP         uint32
	/* Options */
	Data []byte
}

type EthernetFrame struct {
	MACdest      []byte
	MACsource    []byte
	TPID         []byte
	TPIDExtended []byte
	Data         []byte
	// CRC          []byte /* Not considering CRC for now */
}
