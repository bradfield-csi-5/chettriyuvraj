package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/bradfield-csi-5/chettriyuvraj/Prework-5-Traceroute/pkg/traceroute"
	"golang.org/x/sys/unix"
)

/* Default destination */
var DESTADDR unix.SockaddrInet4 = unix.SockaddrInet4{Port: 5431, Addr: [4]byte{0x8E, 0xFA, 0xB7, 0xCE}}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("\n\n No destination found, using default \n\n")
		traceroute.Trace(DESTADDR)
	} else {
		/* Grab and validate first arg which should be dest address */
		args := os.Args[1:]
		destIP := strings.Split(args[0], ".")
		if len(destIP) != 4 {
			log.Fatal("\n\nInvalid IPv4 address \n\n")
		}
		destAddr := unix.SockaddrInet4{Port: 5431, Addr: [4]byte{0x00, 0x00, 0x00, 0x00}}
		for i := 0; i < 4; i++ {
			val, err := strconv.Atoi(destIP[i])
			if err != nil || val > 255 {
				log.Fatal("\n\nError in IPv4 address \n\n")
			}
			destAddr.Addr[i] = byte(val)
		}
		traceroute.Trace(destAddr)
	}

}
