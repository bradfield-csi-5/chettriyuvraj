// TODO: Failing json.decode on running proxy_tester.py

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

const PROXYPORTCLIENT, SERVERPORT = 7543, 9000
const NEWLINE, CARRIAGERETURN = 0x0a, 0x0d
const CACHEPATH = "./cache"
const MAXLISTENQUEUE = 3

var cacheconf CacheConf = CacheConf{
	ProxyCachePath: CACHEPATH,
	Server: []Location{
		{Path: "", ProxyPort: -1},  // always store in order of least to most specific
		{Path: "/", ProxyPort: -1}, // always store in order of least to most specific
		{Path: "/local", ProxyPort: -1},
		{Path: "/proxy", ProxyPath: [4]byte{0x7F, 0x00, 0x00, 0x01}, ProxyPort: 9000},
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
	defer close(clientSocketChannel)
	defer close(errorChannel)

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
		err = unix.Listen(socket, MAXLISTENQUEUE)
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

	defer wg.Done()

	for rwClientSocket := range clientSocketChannel {

		/* Prevent repeated code for closing socket */
		defer func() {
			err := unix.Close(rwClientSocket)
			if err != nil {
				errorChannel <- fmt.Errorf("error while closing client non proxy rw socket %v", err)
			}
		}()

		go func(rwClientSocket int) {
			clientData := make([]byte, 4096)
			var matchLocation Location = Location{Path: "", ProxyPort: -1}

			/* Receive data from client */

			fmt.Println("Receiving from client...")
			_, _, err := unix.Recvfrom(rwClientSocket, clientData, 0)
			if err != nil {
				errorChannel <- fmt.Errorf("error while receiving from client %v", err)
				return
			}

			/* Match location with most specific location in proxy config */

			req, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(clientData)))
			for _, location := range cacheconf.Server {
				fmt.Println(req.URL.Path)
				fmt.Println(location.Path)
				if strings.Contains(req.URL.Path, location.Path) {
					matchLocation = location
				}
			}

			if matchLocation.Path == "" { // return 404
				errorChannel <- fmt.Errorf("location not found on server %v", err)
				return
			}

			/* If location not configured for any proxy - return a simple file from cache folder, this is not a part of caching functionality, just simply reusing func to serve a file */

			if matchLocation.ProxyPort == -1 {
				fmt.Println("No proxy..")
				nonProxyData, err := checkCacheAndServe(CACHEPATH + "/test7")
				if err != nil {
					errorChannel <- fmt.Errorf("error while serving non proxy response %v", err)
					return
				}

				err = unix.Send(rwClientSocket, nonProxyData, 0)
				if err != nil {
					errorChannel <- fmt.Errorf("error while sending non proxy server response to client %v", err)
				}

				return
			}

			/* Check if data exists in cache */

			basePath := path.Base(req.URL.Path)
			filePath := CACHEPATH + "/" + basePath
			cache, err := checkCacheAndServe(filePath)
			if err != nil {
				errorChannel <- fmt.Errorf("error while serving cache %v", err)
				return
			}
			if cache != nil {
				fmt.Println("Serving cache...")
				err = unix.Send(rwClientSocket, cache, 0)
				if err != nil {
					errorChannel <- fmt.Errorf("error while sending  proxy server response to client %v", err)
				}
				return
			}

			/* Data not in cache - get from server, pass to client and cache response */

			serverData, err := proxyToServerAndGetRawResponse(matchLocation, clientData, req)
			if err != nil {
				errorChannel <- fmt.Errorf("error proxying to server and getting raw response %v", err)
				return
			}

			fmt.Println("Passing server response to client...")
			err = unix.Send(rwClientSocket, serverData, 0)
			if err != nil {
				errorChannel <- fmt.Errorf("error while sending server response to client %v", err)
			}

			resp, err := http.ReadResponse(bufio.NewReader(bytes.NewReader(serverData)), req)
			if err != nil {
				errorChannel <- fmt.Errorf("error parsing raw response %v", err)
				return
			}

			cacheResponse(CACHEPATH+"/"+path.Base(req.URL.Path), resp)

		}(rwClientSocket)
	}
}

func handleErrors() {
	defer wg.Done()
	for err := range errorChannel {
		fmt.Printf("non-terminating error %q", err)
	}
}

func checkCacheAndServe(filePath string) ([]byte, error) {
	fmt.Println("Checking cache...")

	serverBuffer := make([]byte, 4096)
	f, err := os.OpenFile(filePath, os.O_RDONLY, 0777)

	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("error while opening cache file %v", err)
	}

	for {
		_, err = f.Read(serverBuffer)
		if err == io.EOF {
			fmt.Println("Done writing cache to buffer...")
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error while writing cache to buffer %v", err)
		}
	}

	response := sampleResponse(string(serverBuffer))
	tempBuffer := bytes.NewBuffer([]byte{})
	err = response.Write(tempBuffer)
	if err != nil {
		return nil, fmt.Errorf("error while writing cache to temp buffer %v", err)
	}

	return tempBuffer.Bytes(), nil

}

func proxyToServerAndGetRawResponse(matchLocation Location, clientData []byte, req *http.Request) ([]byte, error) {

	serverData := make([]byte, 4096)

	fmt.Println("Creating socket to connect to server...")
	serverSocket, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM, 0)
	if err != nil {
		return nil, err
	}

	fmt.Println("Connecting to server...")
	err = unix.Connect(serverSocket, &unix.SockaddrInet4{Port: matchLocation.ProxyPort, Addr: matchLocation.ProxyPath})
	if err != nil {
		return nil, err
	}

	fmt.Println("Passing message to server...")
	err = unix.Send(serverSocket, clientData, 0)
	if err != nil {
		return nil, err
	}

	fmt.Printf("\nGetting response from server...")
	_, _, err = unix.Recvfrom(serverSocket, serverData, 0)
	if err != nil {
		return nil, err
	}

	return serverData, nil
}

func cacheResponse(filePath string, resp *http.Response) error {
	respData := make([]byte, 4096)
	cacheFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
	if err != nil {
		return err
	}

	for {
		_, err := resp.Body.Read(respData)
		fmt.Println(string(respData))
		if err == io.EOF {
			_, err = cacheFile.Write(respData)
			if err != nil {
				return err
			}
			break
		}

		if err != nil {
			return err
		}

		_, err = cacheFile.Write(respData)
		if err != nil {
			return err
		}
	}

	err = cacheFile.Close()
	if err != nil {
		return err
	}

	return nil
}

func sampleResponse(body string) *http.Response {
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
	return t
}
