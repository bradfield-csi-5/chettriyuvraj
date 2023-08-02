// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 221.
//!+

// Netcat1 is a read-only TCP client.
package main

import (
	"io"
	"log"
	"net"
	"os"
	"fmt"
)



func main() {
	ports := []string {"localhost:8000", "localhost:9000", "localhost:10000"}
	goroutinesCreated := false
	for {
		if !goroutinesCreated {
			for _,port := range ports {
				go portTime(port)
			}
		}
		goroutinesCreated = true
	}
}

func portTime(port string) {
	conn, err := net.Dial("tcp", port)
	fmt.Printf("Connected to %s\n", port)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	mustCopy(os.Stdout, conn)
}

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}

