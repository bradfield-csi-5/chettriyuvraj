package main

import (
	"fmt"
	"io"
)

type MyReader struct {
	io.Reader
	str string
	index int
}

func (r *MyReader) Read (p []byte) (n int, err error) {
	if r.index == len(r.str) { // EOF reached
		return 0, io.EOF
	}

	nextIndex := r.index + len(p)
	if r.index + len(p) >= len(r.str) { // if nextIndex exceeds last index of the string - set it to len of string
		nextIndex = len(r.str)
	}
	n = copy(p, r.str[r.index:nextIndex]) // it is legal to copy from a string to a byte slice
	r.index = nextIndex
	
	return n, nil
}

func main() {
	reader := NewReader("Lets try reading this string")
	byteSlice := make([]byte, 5)
	fmt.Println(len(byteSlice))
	for {
		_, err := reader.Read(byteSlice)
		if err == io.EOF {
			fmt.Println("\nEnd of string reached")
			break
		}

		fmt.Printf("%s",byteSlice) // len(p) bytes written each time (until < len(p) bytes remaining to write)

	}
}

func NewReader(s string) *MyReader {
	return &MyReader{str: s}
}