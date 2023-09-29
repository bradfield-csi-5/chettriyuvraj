package ydp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/fnv"
	"time"

	"golang.org/x/sys/unix"
)

type YDPPacket struct {
	ID      uint32 /* Hash32(localaddr, localport, remoteip, remoteport, timestamp) */
	Flags   uint8  /* 7 bits reserved, rightmost bit ack */
	Length  uint16
	Message []byte
}

func SendYDP(message []byte, srcaddr unix.SockaddrInet4, dstaddr unix.SockaddrInet4) (YDPPacket, error) {
	/* Get packet, encode in binary */
	packetYDP, err := NewYDPPacket(message, 0x01, srcaddr, dstaddr)
	if err != nil {
		return YDPPacket{}, err
	}

	encoded, err := packetYDP.encode()
	if err != nil {
		return YDPPacket{}, err
	}

	err = sendUDP(encoded, dstaddr)
	if err != nil {
		return YDPPacket{}, err
	}

	return packetYDP, nil
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

func DecodeYDPPacket(b []byte) (YDPPacket, error) {

	if len(b) < 7 {
		return YDPPacket{}, fmt.Errorf("incomplete YDP headers")
	}

	id := binary.BigEndian.Uint32(b[:4])
	flags := b[4]
	length := binary.BigEndian.Uint16(b[5:7])

	if len(b) != int(length)+7 {
		return YDPPacket{}, fmt.Errorf("message length mismatch in YDP packet")
	}

	return YDPPacket{ID: id, Flags: flags, Length: length, Message: b[7:]}, nil

}
