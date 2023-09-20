package main

import (
	"encoding/binary"
	"fmt"
	"log"

	"golang.org/x/sys/unix"
)

func main() {
	socket, err := unix.Socket(unix.AF_INET, unix.SOCK_DGRAM, 0)
	if err != nil {
		log.Fatalf("Error creating socket %s", err)
	}
	fmt.Println((socket))

	sockaddr := &unix.SockaddrInet4{Port: 53, Addr: [4]byte{0x00, 0x00, 0x00, 0x00}} //Addr: [4]byte{0x6e, 0xe2, 0xb5, 0xb1}}
	ssockaddr := &unix.SockaddrInet4{Port: 53, Addr: [4]byte{0x08, 0x08, 0x08, 0x08}}
	err = unix.Bind(socket, sockaddr)
	if err != nil {
		log.Fatalf("Error binding socket %s", err)
	}

	err = unix.Sendto(socket, dnsMessage(), 0, ssockaddr)
	if err != nil {
		log.Fatalf("Error sending message to DNS server %s", err)
	}

	dnsResponse := make([]byte, 200)
	n, from, err := unix.Recvfrom(socket, dnsResponse, 0)
	if err != nil {
		log.Fatalf("Error recieving message from DNS server %s", err)
	}

	fmt.Println((n))
	fmt.Println((from))
	fmt.Println(dnsResponse)

}

func NewDNSQuery(ID uint16, recursive bool, questions []DNSQuestion) (DNSMessage, error) {
	if len(questions) == 0 {
		return DNSMessage{}, fmt.Errorf("No questions exist")
	}

	return NewDNSMessage(ID, recursive, questions) // increase arguments
}

func NewDNSMessage(ID uint16, recursive bool, questions []DNSQuestion) (DNSMessage, error) { // increase arguments
	header := DNSHeader{ID: ID, Flags: 0x0000, QCount: uint16(len(questions)), AnsCount: 0x0000, RRCount: 0x0000, AddnlRRCount: 0x0000}
	if recursive {
		header.Flags = 0x0080
	}
	return DNSMessage{Header: header, Questions: questions}, nil
}

func (q DNSMessage) executeQuery() []byte {
	socket, err := unix.Socket(unix.AF_INET, unix.SOCK_DGRAM, 0)
	if err != nil {
		log.Fatalf("Error creating socket %s", err)
	}

	sockaddr := &unix.SockaddrInet4{Port: 53, Addr: [4]byte{0x00, 0x00, 0x00, 0x00}}
	ssockaddr := &unix.SockaddrInet4{Port: 53, Addr: [4]byte{0x08, 0x08, 0x08, 0x08}}
	err = unix.Bind(socket, sockaddr)
	if err != nil {
		log.Fatalf("Error binding socket %s", err)
	}

	err = unix.Sendto(socket, q.encode(), 0, ssockaddr)
	if err != nil {
		log.Fatalf("Error sending message to DNS server %s", err)
	}

	dnsResponse := make([]byte, 200)
	_, _, err = unix.Recvfrom(socket, dnsResponse, 0)
	if err != nil {
		log.Fatalf("Error recieving message from DNS server %s", err)
	}

	return dnsResponse
}

func (q DNSMessage) encode() []byte {
	res := []byte{}

	res = binary.BigEndian.AppendUint16(res, q.Header.ID)
	res = binary.BigEndian.AppendUint16(res, q.Header.Flags)
	res = binary.BigEndian.AppendUint16(res, q.Header.QCount)
	res = binary.BigEndian.AppendUint16(res, q.Header.AnsCount)
	res = binary.BigEndian.AppendUint16(res, q.Header.RRCount)
	res = binary.BigEndian.AppendUint16(res, q.Header.AddnlRRCount)

	for _, questions := range q.Questions {
		res = append(res, questions.Name...)
		res = binary.BigEndian.AppendUint16(res, questions.Type)
		res = binary.BigEndian.AppendUint16(res, questions.Class)
	}

	// extend to accomodate remaining fields

	return res
}

func ParseDNSResponse(b []byte) DNSMessage {
	header := DNSHeader{
		ID:           binary.BigEndian.Uint16(b[0:2]),
		Flags:        binary.BigEndian.Uint16(b[2:4]),
		QCount:       binary.BigEndian.Uint16(b[4:6]),
		AnsCount:     binary.BigEndian.Uint16(b[6:8]),
		RRCount:      binary.BigEndian.Uint16(b[8:10]),
		AddnlRRCount: binary.BigEndian.Uint16(b[10:12]),
	}
	questions, _ := ParseDNSQuestions(b, header)
	// answers

	return DNSMessage{Header: header, Questions: questions}
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

func ParseDNSAnswers(b []byte, offset int) ([]DNSAnswer, int) {
	answers := []DNSAnswer{}

	for offset < len(b) {
		answer := DNSAnswer{}

		if b[offset]&0xc0 == 0xc0 {
			labelOffset := int(binary.BigEndian.Uint16(b[offset:offset+2]) & 0x3FFF)
			answer.Name, _ = ParseLabel(b, labelOffset)
		} else {
			answer.Name, offset = ParseLabel(b, offset)
		}
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

func ParseDNSMessageFromEthernet(rawEthernetFrame []byte) DNSMessage {
	ethernetFrame := ParseEthernet(rawEthernetFrame)
	ipDatagram := ParseIPv4(ethernetFrame.Data)
	rawDNS := ipDatagram.Data
	return ParseDNSResponse(rawDNS)
}
