/* Write a program that prints the SHA256 hash of its standard input by default
but supports a command-line flag to print the SHA384 or SHA512 hash instead. */

package main

import (
	"bufio"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"os"
	"strconv"
)

func main() {
	// check for command line flag
	var shaLength int // default
	var shaConvErr error
	for i, arg := range os.Args {
		if arg == "-v" && i != len(os.Args)-1 { // if flag provided and next index exists
			shaLength, shaConvErr = strconv.Atoi(os.Args[i+1])
			break
		}
	}

	if shaLength != 256 && shaLength != 512 && shaLength != 384 { // default
		shaLength = 256
	}

	reader := bufio.NewReader(os.Stdin)
	if line, err := reader.ReadString('\n'); err == nil && shaConvErr == nil {
		fmt.Printf("SHA %d value of %s input is %x", shaLength, line, sha(shaLength, line))
	} else {
		fmt.Println("Error in SHA generation")
	}

}

func sha(shaLength int, data string) []byte {
	switch shaLength {
	case 512:
		return sha512.New().Sum([]byte(data))
	case 384:
		return sha512.New384().Sum([]byte(data))
	default:
		return sha256.New().Sum([]byte(data))
	}

}
