package main

import(
	"fmt"
	"io"
)

type ByteCounter int

// returns only the written byte each time (doesn't even write)
func (l *ByteCounter) Write (p []byte) (n int, err error) {
	return len(p), nil
}

type CountingWriter struct {
	w io.Writer
	Count *int64
}

func (c *CountingWriter) Write(b []byte) (int, error) {
	bytesWritten, err := c.w.Write(b)
	if err != nil {
		return 0, err
	}

	*(c.Count) += int64(bytesWritten)

	return bytesWritten, nil
}


func main(){
	var b *ByteCounter
	b.Write([]byte("Hey hello hi"))

	var cCount int64 = 0
	c := &CountingWriter{w: b, Count: &cCount} // can condense this into a function -> soln to the exercise
	t := []byte("Hello kitty\n Hello bitty")
	c.Write(t)
	fmt.Println(*(c.Count))
	t = []byte("Hi hell")
	c.Write(t)
	fmt.Println(*(c.Count))
}
