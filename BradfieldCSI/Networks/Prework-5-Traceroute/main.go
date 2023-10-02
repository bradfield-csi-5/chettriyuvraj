package main

import (
	"fmt"
	"log"
	"syscall"
)

var CLIENTADDR syscall.SockaddrInet4 = syscall.SockaddrInet4{Port: 5431, Addr: [4]byte{0x7F, 0x00, 0x00, 0x01}}

func main() {

}

/********
- Traceroute
- Decide structure of ICMP and IP requests
- IPv4
	- Version: 4 - 4 bits
	- IHL (Header Length): 20 bytes - 4 bits
	- TOS: 0 - 8 bits
	- Total Length: IP Header length + data length in octets - 16 bits
	- Identification - (any identifier doesn't matter) 16 bits,
	- Flags - 0 (reserved) 1 (Don't fragment) 0 (Last fragment) 3 bits
	- Fragmentation offset (0 since first and only fragment has 0 offset) - 13 bits
	- TTL - inc by 1 each time - 8 bits
	- UL Protocol - 8 bits - 1 for ICMP
	- Header checksum - 16 bits - consider checksum itself as 0 when computing - checksum only on the entire header
	- Source IP - self val - 32 bit
	- Dest IP - dest val - 32 bit

	- Data - ICMP packet

- ICMP (IPv4 payload)
	- Type - 8 for echo message/0 for echo reply - 8 bits
	- Code - 0 for echo reply - 8 bits
	- Checksum - 16 bits header + data
	- Identifier - 16 bits - to help in identifying req and resp echos
	- Sequence no - 16 bits - same as identifier
	- Data - variable - put nothing

- Quickly write NewIPv4Packet and NewTraceRoute funcs which return struct
- Quickly write encode decode methods for two respective structs

********/
