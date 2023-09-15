/*

 */

package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
)

const captureFile, pcapHeaderLength = "./net.cap", 24
const cols = 4

func main() {
	file, err := os.Open(captureFile)
	if err != nil {
		log.Fatalf("error in opening file %v", err)
	}

	header, err := parsePcapHeader(file)
	if err != nil {
		log.Fatalf("Error in parsing pcap header %v", err)
	}

	packetRecords, err := parsePacketRecords(file)
	if err != nil {
		log.Fatalf("Error in parsing packet records %v", err)
	}

	fmt.Println(header)
	fmt.Println(packetRecords)
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
	fmt.Printf("\nTimestamp\tTimestamp2\tCaptured Length\tOriginal Length")
	for _, packetRecord := range packetRecords {
		fmt.Printf("\n%d\t%d\t%d\t%d", packetRecord.Timestamp, packetRecord.Timestamp2, packetRecord.CapturedLen, packetRecord.OriginalLen)
	}
	fmt.Println("Packet len")
	fmt.Println(len(packetRecords))
}

/**
 * Parses ethernet frames from the data extracted from pcap Packet Records
 *
 * @param byte array containing pcap Packet Record payload
 * @return EthernetFrame
 * @return error if found
 *
 * Grabs fields which are there for sure first
 * We consider that TPID(EtherType) is 4 bytes if VLAN else 2 bytes
 * Parses the payload/data of Ethernet frame on the basis of this logic
 *
 * Note: Errors eg. invalid lengths not considered for now
 **/

func ParseEthernetFrame(data []byte) EthernetFrame {
	frame := EthernetFrame{MACdest: data[0:6], MACsource: data[6:12], TPID: data[12:14], TPIDExtended: nil, Data: nil}
	if frame.TPID[0] == 0x81 {
		frame.TPIDExtended = data[14:16]
		frame.Data = data[16:]
	} else {
		frame.Data = data[14:]
	}
	return frame
}
