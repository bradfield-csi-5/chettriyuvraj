package traceroute

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/bits"
)

/* Assuming ICMP packet */
func NewIPv4(SourceIP uint32, DestIP uint32, TTL uint8, Data []byte) IPv4Packet {
	packet := IPv4Packet{
		VersionAndIHL:         0x45,
		TOS:                   0x00,
		TotalLen:              uint16(20 + len(Data)),
		ID:                    0x0000,
		FlagsAndFragmentation: 0x4000, /* don't fragment */
		TTL:                   TTL,
		ULProto:               0x0001, /* ICMP */
		HeaderChecksum:        0x0000,
		SourceIP:              SourceIP,
		DestIP:                DestIP,
	}
	packet.HeaderChecksum = packet.computeHeaderChecksum()
	return packet
}

func NewICMPPacket(Type uint8, Code uint8, Data []byte) ICMPPacket {
	packet := ICMPPacket{Type: Type, Code: Code, ID: 0x0000, SequenceNo: 0x0000, Data: Data}
	return packet
}

/* Assume echo packet and encode id, seq no */
func (p *ICMPPacket) encode() ([]byte, error) {
	buf := new(bytes.Buffer)

	err := binary.Write(buf, binary.BigEndian, p.Type)
	if err != nil {
		return nil, err
	}

	err = binary.Write(buf, binary.BigEndian, p.Code)
	if err != nil {
		return nil, err
	}

	err = binary.Write(buf, binary.BigEndian, p.Checksum)
	if err != nil {
		return nil, err
	}

	err = binary.Write(buf, binary.BigEndian, p.ID)
	if err != nil {
		return nil, err
	}

	err = binary.Write(buf, binary.BigEndian, p.SequenceNo)
	if err != nil {
		return nil, err
	}

	err = binary.Write(buf, binary.BigEndian, p.Data)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func DecodeICMPPacket(b []byte) (ICMPPacket, error) {
	if len(b) < 8 {
		return ICMPPacket{}, fmt.Errorf("invalid icmp packet length")
	}
	return ICMPPacket{Type: b[0], Code: b[1], Checksum: binary.BigEndian.Uint16(b[2:4]), ID: binary.BigEndian.Uint16(b[4:6]), SequenceNo: binary.BigEndian.Uint16(b[6:8]), Data: b[8:]}, nil
}

/* Ignore ID, SequenceNo, always 0 */
func (p *ICMPPacket) computeChecksum() uint16 {
	sum16 := uint16(p.Type)<<8 | uint16(p.Code)
	cur16 := uint16(0x0000)
	b := p.Data

	for i := 0; i < len(b); i += 2 {
		if i+1 < len(p.Data) {
			cur16 = binary.BigEndian.Uint16(b[i : i+2])
		} else {
			cur16 = 0x0000 | (uint16(b[i]) << 8)
		}

		sum32 := uint32(sum16) + uint32(cur16)
		carry16 := uint16(sum32 >> 16)
		sum16 = uint16(sum32) + carry16
	}

	return sum16
}

/* Assume echo packet and encode id, seq no */
func (p *IPv4Packet) encode() ([]byte, error) {
	buf := new(bytes.Buffer)

	err := binary.Write(buf, binary.BigEndian, p.VersionAndIHL)
	if err != nil {
		return nil, err
	}

	err = binary.Write(buf, binary.BigEndian, p.TOS)
	if err != nil {
		return nil, err
	}

	err = binary.Write(buf, binary.BigEndian, p.TotalLen)
	if err != nil {
		return nil, err
	}

	err = binary.Write(buf, binary.BigEndian, p.ID)
	if err != nil {
		return nil, err
	}

	err = binary.Write(buf, binary.BigEndian, p.FlagsAndFragmentation)
	if err != nil {
		return nil, err
	}

	err = binary.Write(buf, binary.BigEndian, p.TTL)
	if err != nil {
		return nil, err
	}

	err = binary.Write(buf, binary.BigEndian, p.ULProto)
	if err != nil {
		return nil, err
	}

	err = binary.Write(buf, binary.BigEndian, p.HeaderChecksum)
	if err != nil {
		return nil, err
	}

	err = binary.Write(buf, binary.BigEndian, p.SourceIP)
	if err != nil {
		return nil, err
	}

	err = binary.Write(buf, binary.BigEndian, p.DestIP)
	if err != nil {
		return nil, err
	}

	err = binary.Write(buf, binary.BigEndian, p.Data)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func DecodeIPv4Packet(b []byte) (IPv4Packet, error) {
	if len(b) < 20 {
		return IPv4Packet{}, fmt.Errorf("invalid ip packet length")
	}
	return IPv4Packet{VersionAndIHL: b[0], TOS: b[1], TotalLen: binary.BigEndian.Uint16(b[2:4]), ID: binary.BigEndian.Uint16(b[4:6]), FlagsAndFragmentation: binary.BigEndian.Uint16(b[6:8]), TTL: b[8], ULProto: b[9], HeaderChecksum: binary.BigEndian.Uint16(b[10:12]), SourceIP: binary.BigEndian.Uint32(b[12:16]), DestIP: binary.BigEndian.Uint32(b[16:20]), Data: b[20:]}, nil
}

func (p *IPv4Packet) computeHeaderChecksum() uint16 {
	sum16 := uint16(p.VersionAndIHL)<<8 | uint16(p.TOS) /* Start with version,IHL + TOS */

	sum32, _ := bits.Add32(uint32(p.TotalLen), uint32(sum16), 0)
	carry16 := uint16(sum32 >> 16)
	sum16 = uint16(sum32) + carry16

	sum32, _ = bits.Add32(uint32(p.ID), uint32(sum16), 0)
	carry16 = uint16(sum32 >> 16)
	sum16 = uint16(sum32) + carry16

	sum32, _ = bits.Add32(uint32(p.FlagsAndFragmentation), uint32(sum16), 0)
	carry16 = uint16(sum32 >> 16)
	sum16 = uint16(sum32) + carry16

	sum32, _ = bits.Add32(uint32(uint16(p.TTL)<<8|uint16(p.ULProto)), uint32(sum16), 0)
	carry16 = uint16(sum32 >> 16)
	sum16 = uint16(sum32) + carry16

	sum32, _ = bits.Add32(uint32(p.HeaderChecksum), uint32(sum16), 0)
	carry16 = uint16(sum32 >> 16)
	sum16 = uint16(sum32) + carry16

	sum32, _ = bits.Add32(uint32(uint16(p.SourceIP)), uint32(sum16), 0)
	carry16 = uint16(sum32 >> 16)
	sum16 = uint16(sum32) + carry16

	sum32, _ = bits.Add32(uint32(uint16(p.SourceIP>>16)), uint32(sum16), 0)
	carry16 = uint16(sum32 >> 16)
	sum16 = uint16(sum32) + carry16

	sum32, _ = bits.Add32(uint32(uint16(p.DestIP)), uint32(sum16), 0)
	carry16 = uint16(sum32 >> 16)
	sum16 = uint16(sum32) + carry16

	sum32, _ = bits.Add32(uint32(uint16(p.DestIP>>16)), uint32(sum16), 0)
	carry16 = uint16(sum32 >> 16)
	sum16 = uint16(sum32) + carry16

	return sum16
}
