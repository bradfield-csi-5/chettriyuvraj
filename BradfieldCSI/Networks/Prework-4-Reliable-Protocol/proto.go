// sender and receiver
// func sender
// create udp socket, bind if you want to, use sendTo to send
// func receive
// create udp socket, bind necessarily, use recvFrom to receive

package main

import (
	"fmt"
	"log"

	"golang.org/x/sys/unix"
)

var PROXYADDR unix.SockaddrInet4 = unix.SockaddrInet4{Port: 64403, Addr: [4]byte{0x7F, 0x00, 0x00, 0x01}}
var RECVSOCKADDR unix.SockaddrInet4 = unix.SockaddrInet4{Port: 5432, Addr: [4]byte{0x7F, 0x00, 0x00, 0x01}}

func main() {

	err := sendUDP([]byte("Hello"), PROXYADDR)
	if err != nil {
		log.Fatal(err)
	}

	recvdmessage, sa, err := receiveUDP(RECVSOCKADDR)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nReceived message %q from socket address %v", recvdmessage, sa)
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
