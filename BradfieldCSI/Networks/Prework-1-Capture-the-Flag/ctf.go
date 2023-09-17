package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
)

const captureFile, pcapHeaderLength, writeFile = "./net.cap", 24, "img.jpg"
const cols = 4

/**
 * Unwrap pcap packets layer by layer and get HTTP image data
 *
 * Parses pcap header
 * Parses packet records
 * Grabs server IP, isolates server packets, sorts them by sequence number, grab data from TCP packets with unique sequence number
 * Write data to file
 **/

func main() {
	readFile, err := os.Open(captureFile)
	if err != nil {
		log.Fatalf("error in opening file %v", err)
	}

	header, err := parsePcapHeader(readFile)
	if err != nil {
		log.Fatalf("Error in parsing pcap header %v", err)
	}
	fmt.Println(header)

	packetRecords, err := parsePacketRecords(readFile)
	if err != nil {
		log.Fatalf("Error in parsing packet records %v", err)
	}

	/* Sort */
	serverIP := ParseIPv4(ParseEthernet(packetRecords[0].Data).Data).DestIP
	serverTCPPackets := filterTCPBySourceIP(packetRecords[0:], serverIP)
	sort.Slice(serverTCPPackets, func(i, j int) bool {
		return serverTCPPackets[i].SequenceNo < serverTCPPackets[j].SequenceNo
	})

	/* Remove dups, grab data and write data */
	hashMap := make(map[uint32]bool)
	httpResponse := []byte{}
	for _, tcpPacket := range serverTCPPackets {
		_, exists := hashMap[tcpPacket.SequenceNo]
		if !exists && tcpPacket.Flags != 0x12 { /* Ignore SYN, ACK packet */
			httpResponse = append(httpResponse, tcpPacket.Data...)
			hashMap[tcpPacket.SequenceNo] = true
		}
	}
	httpResponseBody := ParseHTTP(httpResponse).Body
	writeFile, err := os.OpenFile(writeFile, os.O_RDWR|os.O_CREATE, 0777)
	_, err = writeFile.Write(httpResponseBody)
	if err != nil {
		log.Fatalf("error in opening file to write to %v", err)
	}
}

func parsePcapHeader(file *os.File) (header []byte, err error) {
	header = make([]byte, pcapHeaderLength)

	_, err = file.Read(header)
	if err != nil {
		return nil, fmt.Errorf("error in reading file %v", err)

	}

	return header, nil
	/******
	for i, char := range header {
		if i%cols == 0 {
			fmt.Printf("\n")
		}
		fmt.Printf("%02x", char)
	}
	******/
}

/**
 * Parses packet records from pcap file and returns list of all
 *
 * @param file containing raw binary pcap packet data
 * @return array containing packet records
 * @return error if found
 *
 * Starts from offset 24, since 0-23 are occupied by pcap header
 * For each packet, first parses timestamp, timestamp2, CapturedLen and OriginalLen and stores them array of uint32, converts byte array to uint32 using encoding/binary modules function
 * Using CapturedLen, figures out size of data and then parses it
 * Keeps updating offset at each instance and also checks for EOF incase file is incompelete/EOF has actually been reached
 * Creates a PacketRecord using parsed fields at the end of loop and appends to array
 * Will always return due to EOF at one of the error checks
 *
 * Note: can refactor later (eg grabbing data at once instead of buffering multiple times)
 **/

func parsePacketRecords(file *os.File) (packetRecords []PacketRecord, err error) {
	packetRecords = []PacketRecord{}
	offset := int64(pcapHeaderLength)
	buf := make([]byte, 4)
	fieldBuf := make([]uint32, 4)

	for {
		for i := 0; i < 4; i++ {
			n, err := file.ReadAt(buf, offset)
			if err != nil {
				if err == io.EOF {
					return packetRecords, nil
				}
				return packetRecords, err
			}
			offset += int64(n)
			val := binary.LittleEndian.Uint32(buf)
			fieldBuf[i] = val
		}

		CapturedLen := fieldBuf[2]
		data := make([]byte, CapturedLen)
		n, err := file.ReadAt(data, offset)
		if err != nil {
			if err == io.EOF {
				return packetRecords, nil
			}
			return packetRecords, err
		}
		offset += int64(n)

		packetRecord := PacketRecord{Timestamp: fieldBuf[0], Timestamp2: fieldBuf[1], CapturedLen: fieldBuf[2], OriginalLen: fieldBuf[3], Data: data}
		packetRecords = append(packetRecords, packetRecord)
	}
}

func printPacketRecords(packetRecords []PacketRecord) {
	fmt.Printf("\nIndex\tTimestamp\tTimestamp2\tCaptured Length\tOriginal Length")
	for i, packetRecord := range packetRecords {
		fmt.Printf("\n%d\t%d\t%d\t%d\t%d", i, packetRecord.Timestamp, packetRecord.Timestamp2, packetRecord.CapturedLen, packetRecord.OriginalLen)
	}
}

/**
 * Parses ethernet frames from the data extracted from pcap Packet Records
 *
 * @param byte array containing pcap Packet Record payload
 * @return EthernetFrame
 *
 * Grabs fields which are there for sure first
 * We consider that TPID(EtherType) is 4 bytes if VLAN else 2 bytes
 * Parses the payload/data of Ethernet frame on the basis of this logic
 *
 * Note: Errors eg. invalid lengths not considered for now
 **/

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

/**
 * Parses IPv4 packets from the data extracted from Ethernet frames
 *
 * @param byte array containing Ethernet payload
 * @return IPv4 packet
 *
 * Grabs all fields in order
 * Not considering options field / errors for now
 **/

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
	// packet.Data = data[packet.IHL*4 : packet.TotalLen]
	packet.Data = data[packet.IHL*4:]
	return packet
}

/**
 * Parses TCP packets from the data extracted from IPv4
 *
 * @param byte array containing IPv4
 * @return TCPPacket
 *
 * Grabs all fields in order
 * Not considering options field / errors for now
 **/

func ParseTCP(data []byte) TCPPacket {
	packet := TCPPacket{}
	packet.SourcePort = binary.BigEndian.Uint16(data[0:2])
	packet.DestPort = binary.BigEndian.Uint16(data[2:4])
	packet.SequenceNo = binary.BigEndian.Uint32(data[4:8])
	packet.AckNo = binary.BigEndian.Uint32(data[8:12])
	packet.DataOffset = data[12] & 0xFF >> 4
	packet.Reserved = data[12] & 0x0F
	packet.Flags = data[13] & 0xFF
	packet.WindowSize = binary.BigEndian.Uint16(data[14:16])
	packet.Checksum = binary.BigEndian.Uint16(data[16:18])
	packet.UrgentPtr = binary.BigEndian.Uint16(data[18:20])
	packet.Data = data[packet.DataOffset*4:]
	return packet
}

/**
 * Parses HTTP Messages, both requests and responses
 *
 * @param byte array containing TCP payload
 * @return HTTPMessage
 *
 * Relies on \n\r between parts of HTTP message to parse it
 * Maintains an offset and parses different parts of HTTP message
 **/

func ParseHTTP(data []byte) HTTPMessage {
	request := HTTPMessage{}
	title := []byte{}
	headers := []byte{}
	body := []byte{}
	offset := 0
	newLine := byte(0xA)
	carriageReturn := byte(0xD)
	/* Grab title line */
	for offset = 0; offset < len(data) && (data[offset] != carriageReturn || data[offset+1] != newLine); offset++ {
		title = append(title, data[offset])
	}

	offset += 2
	/* Grab headers */
	for ; offset < len(data) && (data[offset] != carriageReturn || data[offset+1] != newLine || data[offset+2] != carriageReturn || data[offset+3] != newLine); offset++ {
		headers = append(headers, data[offset])
	}

	offset += 4
	/* Grab request body */
	for ; offset < len(data) && (data[offset] != carriageReturn || data[offset+1] != newLine); offset++ {
		body = append(body, data[offset])
	}

	request.Title = title
	request.Headers = headers
	request.Body = body

	return request
}

/**
 * Takes Packet Records, returns TCP Packets filtered by Source IP
 *
 * @param slice of packet records
 * @param source ip
 * @return slice of TCP Packets
 **/

func filterTCPBySourceIP(packetRecords []PacketRecord, sourceIP uint32) []TCPPacket {
	filteredTCPPackets := []TCPPacket{}
	for _, packetRecord := range packetRecords {
		ipPacket := ParseIPv4(ParseEthernet(packetRecord.Data).Data)
		if ipPacket.SourceIP == sourceIP {
			tcpPacket := ParseTCP(ipPacket.Data)
			filteredTCPPackets = append(filteredTCPPackets, tcpPacket)
		}
	}
	return filteredTCPPackets
}
