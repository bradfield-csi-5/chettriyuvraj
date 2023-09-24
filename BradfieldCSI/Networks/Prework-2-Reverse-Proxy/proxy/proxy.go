package main

import (
	"fmt"
	"log"

	"golang.org/x/sys/unix"
)

const PROXYPORTCLIENT, SERVERPORT = 7542, 9000
const NEWLINE, CARRIAGERETURN = 0x0a, 0x0d

var LOOPBACK [4]byte = [4]byte{0x7F, 0x00, 0x00, 0x01}

func main() {
	clientSocket, _, err := ConnectToClient()
	if err != nil {
		log.Fatalf("error connecting to client %v", err)
	}

	err = ForwardClientToServer(clientSocket)
	if err != nil {
		log.Fatalf("error forwarding to server %v", err)
	}
}

func ConnectToClient() (int, unix.Sockaddr, error) {
	fmt.Println("Creating socket to connect to client...")
	socket, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM, 0)
	if err != nil {
		return -1, nil, fmt.Errorf("error creating client socket %v", err)
	}

	sockaddr := &unix.SockaddrInet4{Port: PROXYPORTCLIENT, Addr: LOOPBACK}
	err = unix.Bind(socket, sockaddr)
	if err != nil {
		return -1, nil, fmt.Errorf("error binding socket %v", err)
	}

	fmt.Println("Listening...")
	err = unix.Listen(socket, 1)
	if err != nil {
		return -1, nil, fmt.Errorf("error while listening %v", err)
	}

	nfd, nsa, err := unix.Accept(socket)
	if err != nil {
		return -1, nil, fmt.Errorf("error while accepting %v", err)
	}
	fmt.Println("Accepted...")

	return nfd, nsa, err
}

func ForwardClientToServer(clientSocket int) error {
	clientBuffer := make([]byte, 4096)

	for {
		fmt.Println("Creating socket to connect to server...")
		serverSocket, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM, 0)
		if err != nil {
			return fmt.Errorf("error creating server socket %v", err)
		}

		sockaddr := &unix.SockaddrInet4{Addr: LOOPBACK}
		err = unix.Bind(serverSocket, sockaddr)
		if err != nil {
			return fmt.Errorf("error binding socket %v", err)
		}

		fmt.Println("Receiving from client...")
		n, _, err := unix.Recvfrom(clientSocket, clientBuffer, 0)
		if err != nil {
			log.Fatalf("error while receiving from client %v", err)
		}
		if n == 0 {
			break
		}

		fmt.Println("Connecting to server...")
		err = unix.Connect(serverSocket, &unix.SockaddrInet4{Port: SERVERPORT, Addr: LOOPBACK})
		if err != nil {
			log.Fatalf("error connecting to server %v", err)
		}

		fmt.Printf("\nReceived %q, passing to server...", string(clientBuffer))
		err = unix.Send(serverSocket, clientBuffer, 0)
		if err != nil {
			log.Fatalf("error while sending to server %v", err)
		}

		fmt.Println("\nSuccessfully passed to server, closing connection...")
		err = unix.Close(serverSocket)
		if err != nil {
			log.Fatalf("error while closing server socket %v", err)
		}
	}
	return nil
}
