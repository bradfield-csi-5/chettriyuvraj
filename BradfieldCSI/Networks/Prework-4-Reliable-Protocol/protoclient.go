package main

import (
	"bytes"
	"encoding/binary"
	"hash/fnv"
	"log"
	"time"

	"golang.org/x/sys/unix"
)

type YDPPacket struct {
	ID      uint32
	Flags   uint8 /* 7 bits reserved, rightmost bit ack */
	Length  uint16
	Message []byte
}

var PROXYADDR unix.SockaddrInet4 = unix.SockaddrInet4{Port: 64640, Addr: [4]byte{0x7F, 0x00, 0x00, 0x01}}
var RECVSOCKADDR unix.SockaddrInet4 = unix.SockaddrInet4{Port: 5432, Addr: [4]byte{0x7F, 0x00, 0x00, 0x01}}

func main() {

	err := SendYDP([]byte("Hi guys"), PROXYADDR)
	if err != nil {
		log.Fatal(err)
	}

}

func sendUDP(message []byte, sockaddr unix.SockaddrInet4) error {
	socket, err := unix.Socket(unix.AF_INET, unix.SOCK_DGRAM, 0)
	if err != nil {
		return err
	}

	defer unix.Close(socket)

	err = unix.Sendto(socket, message, 0, &sockaddr)
	if err != nil {
		return err
	}

	return nil
}

func receiveUDP(sockaddr unix.SockaddrInet4) (message []byte, sa unix.Sockaddr, err error) {
	recvdmessage := make([]byte, 4096)

	socket, err := unix.Socket(unix.AF_INET, unix.SOCK_DGRAM, 0)
	if err != nil {
		return nil, nil, err
	}

	defer unix.Close(socket)

	err = unix.Bind(socket, &sockaddr)
	if err != nil {
		return nil, nil, err
	}

	n, sa, err := unix.Recvfrom(socket, recvdmessage, 0)
	if err != nil {
		return nil, nil, err
	}

	return recvdmessage[:n], sa, nil
}

func SendYDP(message []byte, sockaddr unix.SockaddrInet4) error {
	/* Get packet, encode in binary */
	packetYDP, err := NewYDPPacket(message, 0x01, RECVSOCKADDR, sockaddr)
	if err != nil {
		return err
	}

	encoded, err := packetYDP.encode()
	if err != nil {
		return err
	}

	err = sendUDP(encoded, sockaddr)
	if err != nil {
		return err
	}

	return nil
}

func (p *YDPPacket) encode() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.BigEndian, p.ID); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, binary.BigEndian, p.Flags); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, binary.BigEndian, p.Length); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, binary.BigEndian, p.Message); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func NewYDPPacket(message []byte, Flags uint8, sourceaddr unix.SockaddrInet4, destaddr unix.SockaddrInet4) (YDPPacket, error) {
	/* Get unique message id - hash(localaddr, localport, remoteip, remoteport, timestamp)*/

	tsbuf := new(bytes.Buffer)
	err := binary.Write(tsbuf, binary.BigEndian, time.Now().Unix())
	if err != nil {
		return YDPPacket{}, err
	}

	saSlice, err := AddrToByteSlice(sourceaddr)
	if err != nil {
		return YDPPacket{}, err
	}
	dstaSlice, err := AddrToByteSlice(destaddr)
	if err != nil {
		return YDPPacket{}, err
	}

	hashSrc := append(saSlice, dstaSlice...)
	hashSrc = append(hashSrc, tsbuf.Bytes()...)

	hashVal, err := Hash32(hashSrc)
	if err != nil {
		return YDPPacket{}, err
	}

	return YDPPacket{ID: hashVal, Flags: Flags, Length: uint16(len(message)), Message: message}, nil
}

/* Sample YDP packet where timestamp can be specified */
func SampleYDPPacket(message []byte, Flags uint8, sourceaddr unix.SockaddrInet4, destaddr unix.SockaddrInet4, timestamp int64) (YDPPacket, error) {
	/* Get unique message id - hash(localaddr, localport, remoteip, remoteport, timestamp)*/

	tsbuf := new(bytes.Buffer)
	err := binary.Write(tsbuf, binary.BigEndian, timestamp)
	if err != nil {
		return YDPPacket{}, err
	}

	saSlice, err := AddrToByteSlice(sourceaddr)
	if err != nil {
		return YDPPacket{}, err
	}
	dstaSlice, err := AddrToByteSlice(destaddr)
	if err != nil {
		return YDPPacket{}, err
	}

	hashSrc := append(saSlice, dstaSlice...)
	hashSrc = append(hashSrc, tsbuf.Bytes()...)

	hashVal, err := Hash32(hashSrc)
	if err != nil {
		return YDPPacket{}, err
	}

	return YDPPacket{ID: hashVal, Flags: Flags, Length: uint16(len(message)), Message: message}, nil
}

func Hash32(hashSrc []byte) (uint32, error) {
	fnv32 := fnv.New32()
	_, err := fnv32.Write(hashSrc)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint32(fnv32.Sum(nil)), nil
}

func AddrToByteSlice(addr unix.SockaddrInet4) ([]byte, error) {
	buf := new(bytes.Buffer)

	_, err := buf.Write(addr.Addr[:])
	if err != nil {
		return nil, err
	}

	err = binary.Write(buf, binary.BigEndian, uint32(addr.Port))
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// message format
// Header: 32 bit unique message hash, 1 bit ack
// followed by message body

// hash hashSrc
