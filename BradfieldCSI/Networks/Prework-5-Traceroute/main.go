package main

import (
	"github.com/bradfield-csi-5/chettriyuvraj/Prework-5-Traceroute/pkg/traceroute"
	"golang.org/x/sys/unix"
)

var CLIENTADDR unix.SockaddrInet4 = unix.SockaddrInet4{Port: 5431, Addr: [4]byte{0xC0, 0xA8, 0x01, 0x04}}
var DESTADDR unix.SockaddrInet4 = unix.SockaddrInet4{Port: 5431, Addr: [4]byte{0x8E, 0xFA, 0xB7, 0xCE}}

func main() {

	traceroute.Trace(CLIENTADDR, DESTADDR)
}
