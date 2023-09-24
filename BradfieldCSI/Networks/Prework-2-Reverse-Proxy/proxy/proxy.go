package main

import (
	"fmt"
	"log"

	"golang.org/x/sys/unix"
)

const PORT1 = 7542
const PORT2 = 7543
const SERVERPORT = 17542

var ADDR [4]byte = [4]byte{0x7F, 0x00, 0x00, 0x01}

type TCPSocket struct {
	fd int
}

func main() {
	clientSocket, err := NewTCPSocket(PORT1, ADDR)
	if err != nil {
		log.Fatalf("socket error %v", err)
	}

	clientSocketNew, _, err := clientSocket.AcceptConn()
	if err != nil {
		log.Fatalf("socket connect error %v", err)
	}

	serverSocket, err := NewTCPSocket(PORT2, ADDR)
	if err != nil {
		log.Fatalf("socket error %v", err)
	}

	fmt.Println("Connecting to server...")
	err = unix.Connect(serverSocket.fd, &unix.SockaddrInet4{Port: SERVERPORT, Addr: ADDR})
	if err != nil {
		log.Fatalf("error connecting to server %v", err)
	}

	clientBuffer := make([]byte, 4096)
	for {
		n, _, err := unix.Recvfrom(clientSocketNew.fd, clientBuffer, 0)
		if err != nil {
			log.Fatalf("error while receiving from client %v", err)
		}
		if n == 0 {
			break
		}

		fmt.Printf("\nReceived %q, passing to server...", string(clientBuffer))
		err = unix.Send(serverSocket.fd, clientBuffer, 0)
		if err != nil {
			log.Fatalf("error while sending to server %v", err)
		}
		fmt.Printf("\nSuccessfully passed to server...")
	}
}

func NewTCPSocket(PORT int, ADDR [4]byte) (TCPSocket, error) {
	socket, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM, 0)
	if err != nil {
		return TCPSocket{}, fmt.Errorf("error creating socket %v", err)
	}

	sockaddr := &unix.SockaddrInet4{Port: PORT, Addr: ADDR}
	err = unix.Bind(socket, sockaddr)
	if err != nil {
		return TCPSocket{}, fmt.Errorf("error binding socket %v", err)
	}

	return TCPSocket{fd: socket}, nil
}

func (socket *TCPSocket) AcceptConn() (TCPSocket, unix.Sockaddr, error) {
	fmt.Println("Listening...")
	err := unix.Listen(socket.fd, 1)
	if err != nil {
		return TCPSocket{}, nil, fmt.Errorf("error while listening %v", err)
	}

	nfd, nsa, err := unix.Accept(socket.fd)
	if err != nil {
		return TCPSocket{}, nil, fmt.Errorf("error while accepting %v", err)
	}
	fmt.Println("Accepted...")

	return TCPSocket{fd: nfd}, nsa, err
}
