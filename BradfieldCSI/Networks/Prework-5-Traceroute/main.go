package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"time"

	"github.com/bradfield-csi-5/chettriyuvraj/Prework-5-Traceroute/pkg/traceroute"
	"golang.org/x/sys/unix"
)

var CLIENTADDR unix.SockaddrInet4 = unix.SockaddrInet4{Port: 5431, Addr: [4]byte{0xC0, 0xA8, 0x01, 0x04}}
var DESTADDR unix.SockaddrInet4 = unix.SockaddrInet4{Port: 5431, Addr: [4]byte{0x8E, 0xFA, 0xB7, 0xCE}}

func main() {
	icmpPacket := traceroute.NewICMPPacket(0x08, 0x00, []byte{})

	trace(icmpPacket, CLIENTADDR, DESTADDR)
}

func trace(icmpPacket traceroute.ICMPPacket, selfAddr unix.SockaddrInet4, destAddr unix.SockaddrInet4) error {
	recvBuffer := make([]byte, 4096)
	traceMap := make(map[uint16]*traceroute.TraceICMP)

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

			traceMap[counter-1] = &traceroute.TraceICMP{Packet: icmpPacket, StartTime: time.Now()}

			err = unix.Sendto(sendSocket, icmpEncoded, 0, &destAddr)
			if err != nil {
				fmt.Printf("\n\nerror sending message to dest %s", err)
				continue
			}
		}

		for j := 0; j < 3; j++ {
			n, _, err := unix.Recvfrom(recvSocket, recvBuffer, 0)

			if err != nil {
				fmt.Printf("\n\nerror recv message from dest %s", err)
				continue
			}

			recvEncoded := recvBuffer[:n]
			recvIPPacket, err := traceroute.DecodeIPv4Packet(recvEncoded)
			if err != nil {
				fmt.Printf("\n\nerror decoding recvd ipv4 packet %s", err)
				continue
			}

			fmt.Printf("\nIP Data %x", recvIPPacket.Data)
			recvICMPPacket, err := traceroute.DecodeICMPPacket(recvIPPacket.Data)
			if err != nil {
				fmt.Printf("\n\nerror decoding recvd icmp packet %s", err)
				continue
			}

			if recvIPPacket.SourceIP == binary.BigEndian.Uint32(DESTADDR.Addr[:]) && j == 2 { /* Response received from dest and final i.e 3rd packet receieved */
				isComplete = true
			}

			matchingPacket, ok := traceMap[recvICMPPacket.ID]
			if !ok {
				fmt.Printf("Sequence no %d not found in map %s", recvICMPPacket.SequenceNo)
				continue
			}
			matchingPacket.EndTime = time.Now()
			matchingPacket.Response = recvICMPPacket

		}
	}

	return nil
}
