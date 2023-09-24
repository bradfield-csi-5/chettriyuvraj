package main

import (
	"fmt"
	"log"
	"sync"

	"golang.org/x/sys/unix"
)

const PROXYPORTCLIENT, SERVERPORT = 7542, 9000
const NEWLINE, CARRIAGERETURN = 0x0a, 0x0d

// const CACHEPATH = "/cache"

// var cacheconf CacheConf = CacheConf{
// 	ProxyCachePath: CACHEPATH,
// 	Server: []Location{
// 		{Path: "/", ProxyPass: "localhost:7542"},
// 		{Path: "/", ProxyPass: "localhost:7543"},
// 	},
// }

var LOOPBACK [4]byte = [4]byte{0x7F, 0x00, 0x00, 0x01}
var clientSocketChannel chan int = make(chan int, 3)
var errorChannel chan error = make(chan error, 30)
var wg sync.WaitGroup

func main() {

	wg.Add(3)
	go handleErrors()
	go ForwardClientToServer()
	go ConnectToClient()
	wg.Wait()

}

func ConnectToClient() (int, int, error) {

	defer wg.Done()

	fmt.Println("Creating socket to connect to client...")
	socket, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM, 0)
	if err != nil {
		log.Fatalf("error creating client socket %v", err)
	}

	fmt.Println("Binding socket to connect to client...")
	sockaddr := &unix.SockaddrInet4{Port: PROXYPORTCLIENT, Addr: LOOPBACK}
	err = unix.Bind(socket, sockaddr)
	if err != nil {
		log.Fatalf("error binding socket %v", err)
	}

	defer unix.Close(socket)

	for {
		fmt.Println("Listening...")
		err = unix.Listen(socket, 3)
		if err != nil {
			log.Fatalf("error while listening %v", err)
		}

		nfd, _, err := unix.Accept(socket)
		if err != nil {
			errorChannel <- fmt.Errorf("error while accepting %v", err)
			continue
		}

		fmt.Println("Accepted...")
		clientSocketChannel <- nfd
	}

}

func ForwardClientToServer() {
	clientBuffer := make([]byte, 4096)
	serverBuffer := make([]byte, 4096)

	for rwClientSocket := range clientSocketChannel {

		fmt.Println("Creating socket to connect to server...")
		serverSocket, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM, 0)
		if err != nil {
			errorChannel <- fmt.Errorf("error creating server socket %v", err)
			continue
		}

		fmt.Println("Receiving from client...")
		_, _, err = unix.Recvfrom(rwClientSocket, clientBuffer, 0)
		if err != nil {
			errorChannel <- fmt.Errorf("error while receiving from client %v", err)
			continue
		}

		fmt.Println("Connecting to server...")
		err = unix.Connect(serverSocket, &unix.SockaddrInet4{Port: SERVERPORT, Addr: LOOPBACK})
		if err != nil {
			errorChannel <- fmt.Errorf("error connecting to server %v", err)
			continue
		}

		// fmt.Printf("\nReceived %q, passing to server...", string(clientBuffer))
		fmt.Println("Received, passing to server...")
		err = unix.Send(serverSocket, clientBuffer, 0)
		if err != nil {
			errorChannel <- fmt.Errorf("error while sending to server %v", err)
			continue
		}

		fmt.Printf("\nGetting response from server")
		_, _, err = unix.Recvfrom(serverSocket, serverBuffer, 0)
		if err != nil {
			errorChannel <- fmt.Errorf("error while receiving from client %v", err)
			continue
		}

		fmt.Println("\nReceived response from server, closing connection...")
		err = unix.Close(serverSocket)
		if err != nil {
			errorChannel <- fmt.Errorf("error while closing server socket %v", err)
		}

		fmt.Println("Passing server response to client...")
		err = unix.Send(rwClientSocket, serverBuffer, 0)
		if err != nil {
			errorChannel <- fmt.Errorf("error while sending server response to client %v", err)
		}

		err = unix.Close(rwClientSocket)
		if err != nil {
			errorChannel <- fmt.Errorf("error while closing client rw socket %v", err)
		}

	}
}

func handleErrors() {
	defer wg.Done()
	for err := range errorChannel {
		fmt.Printf("non-terminating error %q", err)
	}
}
