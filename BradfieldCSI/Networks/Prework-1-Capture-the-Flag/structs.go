package main

type PacketRecord struct {
	Timestamp   uint32
	Timestamp2  uint32
	CapturedLen uint32
	OriginalLen uint32
	Data        []byte
}

type EthernetFrame struct {
	MACdest      []byte
	MACsource    []byte
	TPID         []byte
	TPIDExtended []byte
	Data         []byte
	// CRC          []byte /* Not considering CRC for now */
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

type TCPPacket struct {
	SourcePort uint16
	DestPort   uint16
	SequenceNo uint32
	AckNo      uint32
	DataOffset uint8
	Reserved   uint8
	Flags      uint8
	WindowSize uint16
	Checksum   uint16
	UrgentPtr  uint16
	/* Options */
	Data []byte
}

type HTTPMessage struct {
	Title   []byte
	Headers []byte
	Body    []byte
}
