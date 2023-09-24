package main

import (
	"fmt"
	"golang.org/x/sys/unix"
	"log"
)

func main() {
	socket, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM, 0)
	if err != nil {
		log.Fatalf("Error creating socket %v", err)
	}

	sockaddr := &unix.SockaddrInet4{Port: 7543, Addr: [4]byte{0x00, 0x00, 0x00, 0x00}}
	err = unix.Bind(socket, sockaddr)
	if err != nil {
		log.Fatalf("Error binding socket %v", err)
	}

	err = unix.Listen(socket, 1)
	if err != nil {
		log.Fatalf("Error while listening %v", err)
	}

	nfd, _, err := unix.Accept(socket)
	if err != nil {
		log.Fatalf("Error while accepting %v", err)
	}

	buffer := make([]byte, 4096)
	for {
		n, _, err := unix.Recvfrom(nfd, buffer, 0)
		if err != nil {
			log.Fatalf("Error while receiving  %v", err)
		}

		if n == 0 {
			break
		}
		fmt.Println(string(buffer))
	}

}
