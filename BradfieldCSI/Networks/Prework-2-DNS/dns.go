package main

import (
	"encoding/binary"
	"fmt"
	"log"

	"golang.org/x/sys/unix"
)

func main() {
	DNSQuery := DNSMessage{
		Header:    DNSHeader{ID: 0x0001, Flags: 0x0100, QCount: 0x0001, AnsCount: 0x0000, RRCount: 0x0000, AddnlRRCount: 0x0000},
		Questions: []DNSQuestion{sampleDNSQuestion()},
	}
	DNSResponse := DNSQuery.send()
	fmt.Printf("\n%x", DNSResponse)
}

func NewDNSMessage(ID uint16, recursive bool, questions []DNSQuestion, answers []DNSAnswer) (DNSMessage, error) { // accomodate remaining arguments
	header := DNSHeader{ID: ID, Flags: 0x0000, QCount: uint16(len(questions)), AnsCount: 0x0000, RRCount: 0x0000, AddnlRRCount: 0x0000}
	if recursive {
		header.Flags = 0x0080
	}
	return DNSMessage{Header: header, Questions: questions}, nil
}

func (q DNSMessage) send() DNSMessage {
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

	rawDNSResponse := make([]byte, 200)
	_, _, err = unix.Recvfrom(socket, rawDNSResponse, 0)
	if err != nil {
		log.Fatalf("Error recieving message from DNS server %s", err)
	}

	return ParseDNSResponse(rawDNSResponse)
}

func (q DNSMessage) encode() []byte {
	res := []byte{}

	res = binary.BigEndian.AppendUint16(res, q.Header.ID)
	res = binary.BigEndian.AppendUint16(res, q.Header.Flags)
	res = binary.BigEndian.AppendUint16(res, q.Header.QCount)
	res = binary.BigEndian.AppendUint16(res, q.Header.AnsCount)
	res = binary.BigEndian.AppendUint16(res, q.Header.RRCount)
	res = binary.BigEndian.AppendUint16(res, q.Header.AddnlRRCount)

	for _, question := range q.Questions {
		res = append(res, question.Name...)
		res = binary.BigEndian.AppendUint16(res, question.Type)
		res = binary.BigEndian.AppendUint16(res, question.Class)
	}

	for _, answer := range q.Answers {
		res = append(res, answer.Name...)
		res = binary.BigEndian.AppendUint16(res, answer.Type)
		res = binary.BigEndian.AppendUint16(res, answer.Class)
		res = binary.BigEndian.AppendUint32(res, answer.TTL)
		res = binary.BigEndian.AppendUint16(res, answer.RDLength)
		res = append(res, answer.RData...)
	}

	// extend to accomodate remaining fields

	return res
}
