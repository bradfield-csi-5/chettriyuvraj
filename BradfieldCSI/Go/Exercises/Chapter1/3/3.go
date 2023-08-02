package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	echo()
	echoOptimized()
}

func echo() {
	s, sep := "", ""
	for i := 1; i < len(os.Args); i++ {
		s += sep + os.Args[i]
		sep = " "
	}
	fmt.Println(s)
}

func echoOptimized() {
	fmt.Println(strings.Join(os.Args[1:], " "))
}
