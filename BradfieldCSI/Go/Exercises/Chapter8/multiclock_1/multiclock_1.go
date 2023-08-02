/* Exercise 8.1: Modify clock2 to accept a port number, and write a program, clockwall, that acts as a client of several 
clock servers at once, reading the times from each one and displaying the results in a table, 
akin to the wall of clocks seen in some business offices. If you have access to geographically distributed computers, 
run instances remotely; otherwise run local instances on different ports with fake time zones.

*/

// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 219.
//!+

// Clock1 is a TCP server that periodically writes the time.
package main

import (
	"io"
	"log"
	"net"
	"time"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("No port number provided")
	}
	portNo := os.Args[1]
	listener, err := net.Listen("tcp", portNo)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err) // e.g., connection aborted
			continue
		}
		go handleConn(conn) // handle concurrent connections
	}
}

func handleConn(c net.Conn) {
	defer c.Close()
	for {
		_, err := io.WriteString(c, c.LocalAddr().String() + " " + time.Now().Format("15:04:05\n"))
		if err != nil {
			return // e.g., client disconnected
		}
		time.Sleep(1 * time.Second)
	}
}
