package main

import (
	"fmt"
)

func main() {
	done := make(chan struct{})
	go func() {
		fmt.Println("performing initialization...")
		<-done
	}()

	done <- struct{}{}
	fmt.Println("initialization done, continuing with rest of program")
}
