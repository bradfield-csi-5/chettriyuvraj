package main

import (
	"encoding/binary"
	"log"
	"os"
	"strings"
	"text/template"

	"golang.org/x/sys/unix"
)

var qtypemap map[uint16]string = map[uint16]string{
	1:  "A",
	2:  "NS",
	5:  "CNAME",
	6:  "SOA",
	15: "MX",
	16: "TXT",
}

var classmap map[uint16]string = map[uint16]string{
	1: "INET",
}

var opcodemap map[uint8]string = map[uint8]string{
	0: "QUERY",
	1: "IQUERY",
	2: "STATUS",
}

var qrmap map[uint8]string = map[uint8]string{
	0: "QUERY",
	1: "REPLY",
}

var rcodemap map[uint8]string = map[uint8]string{
	0: "NOERROR",
	1: "FORMERROR",
	2: "SERVFAIL",
	3: "NXDOMAIN",
}

func main() {
	DNSQuery := DNSMessage{
		Header: DNSHeader{ID: 0x0001, Flags: 0x0000, QCount: 0x0001, AnsCount: 0x0000, RRCount: 0x0000, AddnlRRCount: 0x0000},
		Questions: []DNSQuestion{
			{Name: []byte{0x07, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x03, 0x63, 0x6f, 0x6d, 0x00}, Type: 0x0002, Class: 0x0001},
		},
	}

	DNSResponse := DNSQuery.send()
	DNSResponse.Print()
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

func (m DNSMessage) encode() []byte {
	res := []byte{}

	res = binary.BigEndian.AppendUint16(res, m.Header.ID)
	res = binary.BigEndian.AppendUint16(res, m.Header.Flags)
	res = binary.BigEndian.AppendUint16(res, m.Header.QCount)
	res = binary.BigEndian.AppendUint16(res, m.Header.AnsCount)
	res = binary.BigEndian.AppendUint16(res, m.Header.RRCount)
	res = binary.BigEndian.AppendUint16(res, m.Header.AddnlRRCount)

	for _, question := range m.Questions {
		res = append(res, question.Name...)
		res = binary.BigEndian.AppendUint16(res, question.Type)
		res = binary.BigEndian.AppendUint16(res, question.Class)
	}

	for _, answer := range m.Answers {
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

func (m DNSMessage) Print() {
	// ; <<>> DiG 9.10.6 <<>> example.com
	// ;; global options: +cmd
	// ;; Got answer:
	// ;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 5690
	// ;; flags: qr rd ra; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 1

	// ;; OPT PSEUDOSECTION:
	// ; EDNS: version: 0, flags:; udp: 4096
	// ;; QUESTION SECTION:
	// ;example.com.			IN	A

	// ;; ANSWER SECTION:
	// example.com.		77760	IN	A	93.184.216.34

	// ;; Query time: 7 msec
	// ;; SERVER: fe80::1%15#53(fe80::1%15)
	// ;; WHEN: Thu Sep 21 16:42:07 IST 2023
	// ;; MSG SIZE  rcvd: 56

	templ := `; <<>> YuVi 0.0.0 <<>> {{range .Questions}} {{.Namestr}} {{end}}
	;; Got answer:
	;; ->>HEADER<<- opcode: {{ "OPCODE" | getFlag | OpcodeString }}, status: {{ "RCODE" | getFlag | RcodeString }}, id: {{.Header.ID}}
	;; flags: {{getSetFlagsString}}; QUERY: {{.Header.QCount}}, ANSWER: {{.Header.AnsCount}}, AUTHORITY: {{.Header.RRCount}}, ADDITIONAL: {{.Header.AddnlRRCount}}

	;; QUESTION SECTION:
	{{range	.Questions}}
	;{{.Namestr}}			{{.Type | QtypeString}}	{{.Class | ClassString}}
	{{end}}

	;; ANSWER SECTION:
	{{range .Answers}}
	{{.Namestr}}			{{.Type | QtypeString}}	{{.Class | ClassString}}	{{.TTL}}	{{.RDLength}}	{{.RData}}
	{{end}}

	;; Query time: undefined
	;; SERVER: 8.8.8.8
	;; WHEN: undefined
	;; MSG SIZE  rcvd: undefined` // TODO: Msg size, time, parse RDATA
	templMessage, err := template.New("DNSMessage").Funcs(template.FuncMap{"ClassString": ClassString, "QtypeString": QtypeString, "getFlag": m.getFlag, "OpcodeString": OpcodeString, "QrString": QrString, "RcodeString": RcodeString, "getSetFlagsString": m.getSetFlagsString}).Parse(templ)
	if err != nil {
		log.Fatalf("Error while parsing DNS Message template %v", err)
	}

	err = templMessage.Execute(os.Stdout, m)
	if err != nil {
		log.Fatalf("Error while executing DNS Message template %v", err)
	}
}

func (m DNSMessage) getFlagMap() map[string]uint8 {
	flags := m.Header.Flags
	return map[string]uint8{
		"QR":     uint8((flags & 0b1000000000000000) >> 15),
		"OPCODE": uint8((flags & 0b0111100000000000) >> 11),
		"AA":     uint8((flags & 0b0000010000000000) >> 10),
		"TC":     uint8((flags & 0b0000001000000000) >> 9),
		"RD":     uint8((flags & 0b0000000100000000) >> 8),
		"RA":     uint8((flags & 0b0000000010000000) >> 7),
		"Z":      uint8((flags & 0b0000000001110000) >> 4),
		"RCODE":  uint8((flags & 0b0000000000001111)),
	}
}

func (m DNSMessage) getFlag(f string) uint8 {
	flagMap := m.getFlagMap()
	return flagMap[f]
}

func (m DNSMessage) getSetFlagsString() string {
	flagMap := m.getFlagMap()
	res := []string{}
	for key := range flagMap {
		if key != "OPCODE" && key != "RCODE" && flagMap[key] > 0 {
			res = append(res, key)
		}
	}
	return strings.Join(res, " ")
}

func ClassString(class uint16) string {
	return classmap[class]
}

func QtypeString(qtype uint16) string {
	return qtypemap[qtype]
}

func OpcodeString(opcode uint8) string {
	return opcodemap[opcode]
}

func QrString(qr uint8) string {
	return qrmap[qr]
}

func RcodeString(rcode uint8) string {
	return rcodemap[rcode]
}
