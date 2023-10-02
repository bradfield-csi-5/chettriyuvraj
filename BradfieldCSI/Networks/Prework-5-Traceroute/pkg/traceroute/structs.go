package traceroute

type ICMPPacket struct {
	Type       uint8
	Code       uint8
	Checksum   uint16
	ID         uint16
	SequenceNo uint16
	Data       []byte
}

type IPv4Packet struct {
	VersionAndIHL         uint8
	TOS                   uint8
	TotalLen              uint16
	ID                    uint16
	FlagsAndFragmentation uint16
	TTL                   uint8
	ULProto               uint8
	HeaderChecksum        uint16
	SourceIP              uint32
	DestIP                uint32
	Data                  []byte
}

/* - IPv4
- Version: 4 - 4 bits
- IHL (Header Length): 20 bytes - 4 bits
- TOS: 0 - 8 bits
- Total Length: IP Header length + data length in octets - 16 bits
- Identification - (any identifier doesn't matter) 16 bits,
- Flags - 0 (reserved) 1 (Don't fragment) 0 (Last fragment) 3 bits
- Fragmentation offset (0 since first and only fragment has 0 offset) - 13 bits
- TTL - inc by 1 each time - 8 bits
- UL Protocol - 8 bits - 1 for ICMP
- Header checksum - 16 bits - consider checksum itself as 0 when computing - checksum only on the entire header
- Source IP - self val - 32 bit
- Dest IP - dest val - 32 bit

- Data - ICMP packet */
