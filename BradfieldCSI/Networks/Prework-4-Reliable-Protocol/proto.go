package main

import (
	"log"

	"github.com/bradfield-csi-5/chettriyuvraj/Prework-4-Reliable-Protocol/pkg/ydp"
	"golang.org/x/sys/unix"
)

var PROXYADDR unix.SockaddrInet4 = unix.SockaddrInet4{Port: 52372, Addr: [4]byte{0x7F, 0x00, 0x00, 0x01}}
var RECVSOCKADDR unix.SockaddrInet4 = unix.SockaddrInet4{Port: 5432, Addr: [4]byte{0x7F, 0x00, 0x00, 0x01}}

func main() {

	err := ydp.SendYDP([]byte("Hi guys"), RECVSOCKADDR, PROXYADDR)
	if err != nil {
		log.Fatal(err)
	}

}
