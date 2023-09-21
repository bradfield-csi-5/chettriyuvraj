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
		question.Namestr = LabelToString(question.Name)
		offset += 4

		questions = append(questions, question)
	}

	return questions, offset
}

func ParseDNSAnswers(b []byte, header DNSHeader, offset int) ([]DNSAnswer, int) {
	answers := []DNSAnswer{}

	for i := 0; i < int(header.AnsCount); i++ {
		answer := DNSAnswer{}
		if b[offset]&0xc0 == 0xc0 {
			labelOffset := ConvPointerToOffset(b, offset)
			name, _ := ParseLabel(b, labelOffset)
			answer.Name = b[offset : offset+2]
			answer.Namestr = LabelToString(name)
		} else {
			answer.Name, offset = ParseLabel(b, offset)
			answer.Namestr = LabelToString(answer.Name)
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

func ConvPointerToOffset(b []byte, offset int) int {
	labelOffset := int(binary.BigEndian.Uint16(b[offset:offset+2]) & 0x3FFF)
	return labelOffset
}

func LabelToString(label []byte) string {
	res := []byte{}
	for offset := 0; label[offset] != 0x00; offset += int(label[offset]) + 1 {
		if len(res) > 0 {
			res = append(res, 0x2E)
		}
		res = append(res, label[offset+1:offset+int(label[offset])+1]...)
	}
	return string(res)
}
