package main

import (
	"fmt"
	"bufio"
	"bytes"
)

type LineCounter int
type WordCounter int

// counting only the written line each time instead of adding to LineCounter value
func (l *LineCounter) Write (p []byte) (n int, err error) {
	scanner := bufio.NewScanner(bytes.NewReader(p)) // new scanner from new reader from byte slice
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}
	return lineCount, nil
}

// counting only the written line each time instead of adding to LineCounter value
func (w *WordCounter) Write (p []byte) (n int, err error) {
	scanner := bufio.NewScanner(bytes.NewReader(p)) // new scanner from new reader from byte slice
	scanner.Split(bufio.ScanWords)
	wordCount := 0
	for scanner.Scan() {
		wordCount++
	}
	return wordCount, nil
}

func main () {
	var l LineCounter
	var w WordCounter
	text := []byte("Hello kitty\n Hello bitty")
	// l.Write(text)
	// w.Write(text)/
	fmt.Println(l.Write(text))
	fmt.Println(w.Write(text))
}

