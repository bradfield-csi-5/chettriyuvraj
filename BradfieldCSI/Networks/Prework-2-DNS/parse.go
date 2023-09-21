package main

import (
	"encoding/binary"
)

func ParseDNSResponse(b []byte) DNSMessage {
	header := DNSHeader{
		ID:           binary.BigEndian.Uint16(b[0:2]),
		Flags:        binary.BigEndian.Uint16(b[2:4]),
		QCount:       binary.BigEndian.Uint16(b[4:6]),
		AnsCount:     binary.BigEndian.Uint16(b[6:8]),
		RRCount:      binary.BigEndian.Uint16(b[8:10]),
		AddnlRRCount: binary.BigEndian.Uint16(b[10:12]),
	}
	questions, offset := ParseDNSQuestions(b, header)
	answers, offset := ParseDNSAnswers(b, header, offset)
	return DNSMessage{Header: header, Questions: questions, Answers: answers}
}

func ParseDNSQuestions(b []byte, header DNSHeader) ([]DNSQuestion, int) {
	questions := []DNSQuestion{}
	offset := 12

	for i := uint16(0); i < header.QCount; i++ {
		question := DNSQuestion{}

		question.Name, offset = ParseLabel(b, offset)
		question.Type = binary.BigEndian.Uint16(b[offset : offset+2])
		question.Class = binary.BigEndian.Uint16(b[offset+2 : offset+4])
		offset += 4

		questions = append(questions, question)
	}

	return questions, offset
}

func ParseDNSAnswers(b []byte, header DNSHeader, offset int) ([]DNSAnswer, int) {
	answers := []DNSAnswer{}

	for i := 0; i < int(header.AnsCount); i++ {
		answer := DNSAnswer{}
		answer.Name = b[offset : offset+2]
		answer.Type = binary.BigEndian.Uint16(b[offset+2 : offset+4])
		answer.Class = binary.BigEndian.Uint16(b[offset+4 : offset+6])
		answer.TTL = binary.BigEndian.Uint32(b[offset+6 : offset+10])
		answer.RDLength = binary.BigEndian.Uint16(b[offset+10 : offset+12])
		answer.RData = b[offset+12 : offset+12+int(answer.RDLength)]
		answers = append(answers, answer)
		offset += 12 + int(answer.RDLength)
	}

	return answers, offset
}

func ParseLabel(b []byte, offset int) ([]byte, int) {
	res := []byte{}

	for ; b[offset] != 0x00; offset += int(b[offset]) + 1 {
		res = append(res, b[offset:offset+int(b[offset])+1]...)
	}
	res = append(res, 0x00)
	offset++

	return res, offset
}

func ConvPointerToLabel(b []byte, offset int) ([]byte, int) {
	/*****
		if b[offset]&0xc0 == 0xc0 {
		labelOffset := int(binary.BigEndian.Uint16(b[offset:offset+2]) & 0x3FFF)
		answer.Name, _ = ParseLabel(b, labelOffset)
		answer.Name, _ = ConvPointerToLabel(b, offset)
	} ******/
	labelOffset := int(binary.BigEndian.Uint16(b[offset:offset+2]) & 0x3FFF)
	return ParseLabel(b, labelOffset)
}

func ParseIPv4(data []byte) IPv4Packet {
	packet := IPv4Packet{}
	packet.Version = data[0] & 0xF0 >> 4                                /* bits 0-3 */
	packet.IHL = data[0] & 0x0F                                         /* bits 4-7 */
	packet.DSCP = (data[1] & 0xFC) >> 2                                 /* bits 8-13 */
	packet.ECN = data[1] & 0x03                                         /* bits 14-15 */
	packet.TotalLen = binary.BigEndian.Uint16(data[2:4])                /* bits 16-31 bytes 2,3 */ // (uint16(data[2]) & 0xFF00) & (0x00FF & uint16(data[3]))
	packet.Id = binary.BigEndian.Uint16(data[4:6])                      /* bytes 4,5 */
	packet.Flags = (data[6] & 0xE0) >> 5                                /* byte 6 bits 0-3 */
	packet.FragmentOffset = binary.BigEndian.Uint16(data[6:8]) & 0x1FFF /* byte 6, 7, last 13 bits */
	packet.TTL = data[8]                                                /* byte 8 */
	packet.Protocol = data[9]                                           /* byte 9 */
	packet.HeaderChecksum = binary.BigEndian.Uint16(data[10:12])        /* bytes 10,11 */
	packet.SourceIP = binary.BigEndian.Uint32(data[12:16])              /* bytes 12-15 */
	packet.DestIP = binary.BigEndian.Uint32(data[16:20])                /* bytes 16-19 */
	packet.Data = data[packet.IHL*4:]
	return packet
}

func ParseEthernet(data []byte) EthernetFrame {
	frame := EthernetFrame{MACdest: data[0:6], MACsource: data[6:12], TPID: data[12:14], TPIDExtended: nil, Data: nil}
	if frame.TPID[0] == 0x81 {
		frame.TPIDExtended = data[14:16]
		frame.Data = data[16:]
	} else {
		frame.Data = data[14:]
	}
	return frame
}

func ParseUDP(data []byte) UDPPacket {
	packet := UDPPacket{}
	packet.SourcePort = binary.BigEndian.Uint16(data[0:2])
	packet.DestPort = binary.BigEndian.Uint16(data[2:4])
	packet.Length = binary.BigEndian.Uint16(data[4:6])
	packet.Checksum = binary.BigEndian.Uint16(data[6:8])
	packet.Data = data[8:]
	return packet
}
