package traceroute

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math/bits"
	"os"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"golang.org/x/sys/unix"
)

func Trace(destAddr unix.SockaddrInet4) error {
	icmpPacket := NewICMPPacket(0x08, 0x00, []byte{})
	recvBuffer := make([]byte, 4096)
	traceMap := make(map[uint16]*TraceICMP)

	fmt.Printf("\n\nxx----------------UGLY PRINT WARNING----------------------xx\n\n\n")

	fmt.Printf("\n\nxx----------------TRACING IS VERY FAST, SLEEP DELIBERATELY ADDED----------------------xx\n\n\n")
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Start IP", "Route IP", "RTT"})

	t.Render()
	t.ResetHeaders()

	/* Sending socket  */
	sendSocket, err := unix.Socket(unix.AF_INET, unix.SOCK_RAW, unix.IPPROTO_ICMP)
	if err != nil {
		log.Fatalf("error creating socket for server %s", err)
	}
	defer unix.Close(sendSocket)

	/* Recv socket with timeout */
	recvSocket, err := unix.Socket(unix.AF_INET, unix.SOCK_RAW, unix.IPPROTO_ICMP)
	if err != nil {
		log.Fatalf("error creating recv socket for server %s", err)
	}
	defer unix.Close(recvSocket)

	tv := unix.Timeval{Sec: int64(10), Usec: 0}
	err = unix.SetsockoptTimeval(recvSocket, unix.SOL_SOCKET, unix.SO_RCVTIMEO, &tv)
	if err != nil {
		log.Fatalf("error setting recv sock opt %s", err)
	}

	/* Traceroute */
	counter := uint16(0)
	isComplete := false
	for i := 1; isComplete == false; i++ {

		err = unix.SetsockoptInt(sendSocket, unix.IPPROTO_IP, unix.IP_TTL, i)
		if err != nil {
			log.Fatalf("error setting sock opt %s", err)
		}

		for j := 0; j < 3; j++ {
			/* Add ID, Sequence no and checksum */
			icmpPacket.ID = counter
			icmpPacket.SequenceNo = counter
			icmpPacket.Checksum = icmpPacket.ComputeChecksum()
			counter += 1

			icmpEncoded, err := icmpPacket.Encode()
			if err != nil {
				fmt.Printf("\n\nerror encoding message to dest %s", err)
				continue
			}

			traceMap[icmpPacket.ID] = &TraceICMP{Packet: icmpPacket, StartTime: time.Now()}

			err = unix.Sendto(sendSocket, icmpEncoded, 0, &destAddr)
			if err != nil {
				fmt.Printf("\n\nerror sending message to dest %s", err)
				continue
			}
		}

		for j := 0; j < 3; j++ {
			n, _, err := unix.Recvfrom(recvSocket, recvBuffer, 0)

			if err != nil {
				// fmt.Printf("\n\nerror recv message from dest %s", err)
				t.AppendRow([]interface{}{"-----------", err, "-----------"})
				continue
			}

			recvEncoded := recvBuffer[:n]
			recvIPPacket, err := DecodeIPv4Packet(recvEncoded)
			if err != nil {
				// fmt.Printf("\n\nerror decoding recvd ipv4 packet %s", err)
				t.AppendRow([]interface{}{"-----------", err, "-----------"})
				continue
			}

			recvICMPPacket, err := DecodeICMPPacket(recvIPPacket.Data)
			if err != nil {
				// fmt.Printf("\n\nerror decoding recvd icmp packet %s", err)
				t.AppendRow([]interface{}{"-----------", err, "-----------"})
				continue
			}

			if recvIPPacket.SourceIP == binary.BigEndian.Uint32(destAddr.Addr[:]) && j == 2 { /* Response received from dest and final i.e 3rd packet receieved */
				isComplete = true
			}

			matchingPacket, ok := traceMap[recvICMPPacket.ID]
			if !ok {
				// fmt.Printf("\n\nSequence no %d not found", recvICMPPacket.SequenceNo)
				t.AppendRow([]interface{}{"-----------", err, "-----------"})
				continue
			}
			matchingPacket.EndTime = time.Now()
			matchingPacket.Response = recvICMPPacket

			t.AppendRow([]interface{}{matchingPacket.Packet.ID, uint32toIPv4(recvIPPacket.DestIP), uint32toIPv4(recvIPPacket.SourceIP), matchingPacket.EndTime.Sub(matchingPacket.StartTime)})

		}
		time.Sleep(time.Second * 3)
		t.AppendSeparator()
		t.Render()
		t.ResetRows()
	}

	return nil
}
func NewICMPPacket(ptype uint8, code uint8, data []byte) ICMPPacket {
	packet := ICMPPacket{Type: ptype, Code: code, ID: 0x0000, SequenceNo: 0x0000, Data: data}
	packet.Checksum = packet.ComputeChecksum()
	return packet
}

/* Assume echo packet and encode id, seq no */
func (p *ICMPPacket) Encode() ([]byte, error) {
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
	/* Check type */
	packetType := b[0]
	if packetType == 0x00 || packetType == 0x08 {
		return DecodeICMPEchoPacket(b)
	} else {
		return DecodeICMPOtherPacket(b)
	}
}

func DecodeICMPEchoPacket(b []byte) (ICMPPacket, error) {
	if len(b) < 8 {
		return ICMPPacket{}, fmt.Errorf("invalid icmp packet length")
	}
	return ICMPPacket{Type: b[0], Code: b[1], Checksum: binary.BigEndian.Uint16(b[2:4]), ID: binary.BigEndian.Uint16(b[4:6]), SequenceNo: binary.BigEndian.Uint16(b[6:8]), Data: b[8:]}, nil
}

/* Packets other than echo resp/request that are possible i.e all types other than timestamp/information which have this same format */
func DecodeICMPOtherPacket(b []byte) (ICMPPacket, error) {
	if len(b) < 8 {
		return ICMPPacket{}, fmt.Errorf("invalid icmp ttl exceeded packet length")
	}

	ipPacket, err := DecodeIPv4Packet(b[8:])
	if err != nil {
		return ICMPPacket{}, fmt.Errorf("invalid ipv4 packet in ttl exceeded packet")
	}

	icmppacket, err := DecodeICMPPacket(ipPacket.Data)
	if err != nil {
		return ICMPPacket{}, fmt.Errorf("invalid icmp packet in ttl exceeded packet")
	}

	return icmppacket, nil
}

func (p *ICMPPacket) ComputeChecksum() uint16 {
	sum16 := uint16(p.Type)<<8 | uint16(p.Code)
	sum16 += p.SequenceNo + p.ID
	cur16 := uint16(0x0000)
	b := p.Data

	for i := 0; i < len(b); i += 2 {
		if i+1 < len(p.Data) {
			cur16 = binary.BigEndian.Uint16(b[i : i+2])
		} else {
			cur16 = 0x0000 | ((uint16(b[i])) << 8)
		}

		sum32 := uint32(sum16) + uint32(cur16)
		carry16 := uint16(sum32 >> 16)
		sum16 = uint16(sum32) + carry16
	}

	return 0xFFFF ^ sum16
}

func DecodeIPv4Packet(b []byte) (IPv4Packet, error) {
	if len(b) < 20 {
		return IPv4Packet{}, fmt.Errorf("invalid ip packet length")
	}
	/* Arbitrarily assuming packet as 20 bytes */
	return IPv4Packet{VersionAndIHL: b[0], TOS: b[1], TotalLen: binary.BigEndian.Uint16(b[2:4]), ID: binary.BigEndian.Uint16(b[4:6]), FlagsAndFragmentation: binary.BigEndian.Uint16(b[6:8]), TTL: b[8], ULProto: b[9], HeaderChecksum: binary.BigEndian.Uint16(b[10:12]), SourceIP: binary.BigEndian.Uint32(b[12:16]), DestIP: binary.BigEndian.Uint32(b[16:20]), Data: b[20:]}, nil
}

func uint32toIPv4(val uint32) (res [4]byte) {
	for i := 3; i >= 0; i-- {
		res[i] = uint8(0x00000000 | val)
		val >>= 8
	}
	return res
}

/**
 * Unused functions
 * Had assumed I would have to frame my own IP packets entirely to implement tracing with SOCK_RAW
 * Turned out not to be the case
 **/

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

/* Assuming ICMP packet */
func NewIPv4(sourceIP uint32, destIP uint32, ttl uint8, data []byte) IPv4Packet {

	packet := IPv4Packet{
		VersionAndIHL:         0x45,
		TOS:                   0x00,
		TotalLen:              0x14 + uint16(len(data)),
		ID:                    0x0000,
		FlagsAndFragmentation: 0x4000, /* don't fragment */
		TTL:                   ttl,
		ULProto:               0x0001, /* ICMP */
		HeaderChecksum:        0x0000,
		SourceIP:              sourceIP,
		DestIP:                destIP,
	}
	packet.HeaderChecksum = packet.computeHeaderChecksum()
	return packet
}

func (p *IPv4Packet) Encode() ([]byte, error) {
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
