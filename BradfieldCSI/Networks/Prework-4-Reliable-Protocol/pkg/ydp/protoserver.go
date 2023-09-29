package ydp

import (
	"golang.org/x/sys/unix"
)

type YDPServer struct {
}

func (server *YDPServer) RecvAndAck(selfaddr unix.SockaddrInet4, dstaddr unix.SockaddrInet4, socket int) error {
	data := make([]byte, 4096)

	n, _, err := unix.Recvfrom(socket, data, 0)
	if err != nil {
		return err
	}

	packet, err := DecodeYDPPacket(data[:n])
	if err != nil {
		return err
	}

	ackpacket := YDPPacket{ID: packet.ID, Flags: 0x01, Length: 6, Message: []byte("server")}

	encoded, err := ackpacket.encode()
	if err != nil {
		return err
	}

	err = sendUDP(encoded, dstaddr)
	if err != nil {
		return err
	}

	return nil
}
