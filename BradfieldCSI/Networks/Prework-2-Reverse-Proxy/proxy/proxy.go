package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"

	"golang.org/x/sys/unix"
)

const PROXYPORTCLIENT, SERVERPORT = 7542, 9000
const NEWLINE, CARRIAGERETURN = 0x0a, 0x0d

const CACHEPATH = "./cache"

var cacheconf CacheConf = CacheConf{
	ProxyCachePath: CACHEPATH,
	Server: []Location{
		{Path: "/proxy", ProxyPath: [4]byte{0x7F, 0x00, 0x00, 0x01}, ProxyPort: 9000},
		{Path: "/local", ProxyPort: -1}, // store in order of least to most specific
	},
}

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

		var matchLocation Location = Location{Path: "", ProxyPort: -1}

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

		req, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(clientBuffer)))
		for _, location := range cacheconf.Server {
			if strings.Contains(req.URL.Path, location.Path) { // location.Path stored in order of least to most specific
				fmt.Println("match")
				fmt.Println(location)
				matchLocation = location
			}
		}

		fmt.Println(matchLocation)

		if matchLocation.Path == "" { // return 404
			errorChannel <- fmt.Errorf("location not found on server %v", err)
			continue
		}

		if matchLocation.ProxyPort == -1 { // no proxy
			fmt.Println("No proxy..")
			response := sampleResponse("no proxy")
			tempBuffer := bytes.NewBuffer(serverBuffer)
			err = response.Write(tempBuffer)
			if err != nil {
				errorChannel <- fmt.Errorf("error while writing non proxy response to buffer %v", err)
				continue
			}

			err = unix.Send(rwClientSocket, tempBuffer.Bytes(), 0)
			if err != nil {
				errorChannel <- fmt.Errorf("error while sending non proxy server response to client %v", err)
			}

			err = unix.Close(rwClientSocket)
			if err != nil {
				errorChannel <- fmt.Errorf("error while closing client non proxy rw socket %v", err)
			}
			continue
		}

		fmt.Println("Checking cache...")
		basePath := path.Base(req.URL.Path)
		filePath := CACHEPATH + "/" + basePath
		fmt.Println(filePath)
		f, err := os.OpenFile(filePath, os.O_RDONLY, 0777)
		if err == nil {
			fmt.Println("Serving cache...")
			for {
				_, err = f.Read(serverBuffer)
				if err == io.EOF {
					fmt.Println("Done writing cache to buffer...")
					break
				}
				if err != nil {
					errorChannel <- fmt.Errorf("error while writing cache to buffer %v", err)
					break
				}
			}
			if err == io.EOF || err == nil {

				response := sampleResponse(string(serverBuffer))
				tempBuffer := bytes.NewBuffer([]byte{})
				err = response.Write(tempBuffer)

				err = unix.Send(rwClientSocket, tempBuffer.Bytes(), 0)
				if err != nil {
					errorChannel <- fmt.Errorf("error while sending  proxy server response to client %v", err)
				}

				err = unix.Close(rwClientSocket)
				if err != nil {
					errorChannel <- fmt.Errorf("error while closing client  proxy rw socket %v", err)
				}
			}
			continue
		}

		fmt.Println("File not in cache")
		fmt.Println("Connecting to server...")
		err = unix.Connect(serverSocket, &unix.SockaddrInet4{Port: matchLocation.ProxyPort, Addr: matchLocation.ProxyPath})
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

func sampleResponse(body string) *http.Response { // Giving error when sending to curl
	t := &http.Response{
		StatusCode:    http.StatusOK,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Body:          ioutil.NopCloser(bytes.NewBufferString(body)),
		Request:       &http.Request{Method: "GET"},
		ContentLength: int64(len(body)),
		Header:        http.Header{},
	}
	t.Header.Set("Content-Length", string(len(body)))        // Set the Date header
	t.Header.Set("Content-Type", "application/octet-stream") // Set the Content-Type header
	t.Header.Set("Server", "My HTTP Server")                 // Set the Server header
	t.Header.Set("Date", "Sun, 24 Sep 2023 14:00:00 GMT")    // Set the Date header
	return t
}

// whenever a request comes in:
// match it with the locations most specific
// check if the file exists in cachePath directory, if yes return it back
// otherwise route to server, save the data when it returns from the server to cache folder
