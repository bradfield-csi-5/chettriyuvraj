package traceroute

import "time"

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

type TraceICMP struct {
	StartTime time.Time
	EndTime   time.Time
	Packet    ICMPPacket
	Response  ICMPPacket
}
