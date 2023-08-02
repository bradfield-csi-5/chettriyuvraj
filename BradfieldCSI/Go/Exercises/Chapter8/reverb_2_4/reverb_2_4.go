/* Modify the reverb2 server to use a sync.WaitGroup per connection to count the number of active echo goroutines. 
When it falls to zero, close the write half of the TCP connection as described in Exercise 8.3. Verify that your 
modified netcat3 client from that exercise waits for the final echoes of multiple concurrent shouts, even after the 
standard input has been closed.
*/


// Note: When EOF is encountered and input.Scan() is released, only then can wg.Wait() start - 
// Solve the question with that understanding in mind


// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 224.

// Reverb2 is a TCP server that simulates an echo.
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
	"sync"
)

func echo(c net.Conn, shout string, delay time.Duration) {
	fmt.Fprintf(c, "Echoing %v\n", shout)
	fmt.Fprintln(c, "\t", strings.ToUpper(shout))
	fmt.Fprintln(c, "\t", strings.ToUpper(shout))
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", shout)
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", strings.ToLower(shout))
}

//!+
func handleConn(c net.Conn) {
	var wg sync.WaitGroup
	input := bufio.NewScanner(c)


	for input.Scan() {
		wg.Add(1)
		go func () {
			echo(c, input.Text(), 1*time.Second)
			wg.Done()
		}()
	}

	fmt.Println("Waiting!")
	wg.Wait()
	fmt.Println("Closing Write Section!")
	c.(*net.TCPConn).CloseWrite()
}

//!-

func main() {
	l, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Print(err) // e.g., connection aborted
			continue
		}
		go handleConn(conn)
	}

}