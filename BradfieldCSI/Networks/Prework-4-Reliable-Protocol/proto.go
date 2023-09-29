package main

import (
	"log"
	"strconv"
	"sync"

	"github.com/bradfield-csi-5/chettriyuvraj/Prework-4-Reliable-Protocol/pkg/ydp"
	"golang.org/x/sys/unix"
)

var PROXYADDR unix.SockaddrInet4 = unix.SockaddrInet4{Port: 54411, Addr: [4]byte{0x7F, 0x00, 0x00, 0x01}}
var CLIENTADDR unix.SockaddrInet4 = unix.SockaddrInet4{Port: 5431, Addr: [4]byte{0x7F, 0x00, 0x00, 0x01}}
var SERVERADDR unix.SockaddrInet4 = unix.SockaddrInet4{Port: 1234, Addr: [4]byte{0x7F, 0x00, 0x00, 0x01}}
var PROXYADDR2 unix.SockaddrInet4 = unix.SockaddrInet4{Port: 54412, Addr: [4]byte{0x7F, 0x00, 0x00, 0x01}}

var wg sync.WaitGroup

func main() {

	wg.Add(1)
	go server()
	wg.Add(1)
	go client()
	wg.Wait()

}

func server() {
	defer wg.Done()

	server := ydp.YDPServer{}

	socket, err := unix.Socket(unix.AF_INET, unix.SOCK_DGRAM, 0)
	if err != nil {
		log.Fatalf("error creating socket for server")
	}
	defer unix.Close(socket)

	err = unix.Bind(socket, &SERVERADDR)
	if err != nil {
		log.Fatalf("error binding socket for server")
	}

	for {
		err = server.RecvAndAck(SERVERADDR, PROXYADDR2, socket)
	}
}

func client() {
	defer wg.Done()

	client := ydp.YDPClient{}

	socket, err := unix.Socket(unix.AF_INET, unix.SOCK_DGRAM, 0)
	if err != nil {
		log.Fatalf("error creating socket for client")
	}
	defer unix.Close(socket)

	_, err = unix.FcntlInt(uintptr(socket), unix.F_SETFL, unix.O_NONBLOCK)
	if err != nil {
		log.Fatalf("error setting client socket to non blocking")
	}

	err = unix.Bind(socket, &CLIENTADDR)
	if err != nil {
		log.Fatalf("error binding socket for client")
	}

	for i := 0; i < 5; i++ {

		message := append([]byte("Client "), strconv.Itoa(i)...)

		err = client.Send(message, CLIENTADDR, PROXYADDR, socket)
		if err != nil {
			log.Fatalf("Failed to send ydp packet for packet no %d", i)
			continue
		}
	}

}
