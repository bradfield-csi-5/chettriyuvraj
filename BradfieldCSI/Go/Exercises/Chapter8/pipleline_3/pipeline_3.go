// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 231.

// Pipeline3 demonstrates a finite 3-stage pipeline
// with range, close, and unidirectional channel types.
package main

import "fmt"

//!+
func counter(out chan<- int) {
	fmt.Println("counter")
	for x := 0; x < 100; x++ {
		fmt.Printf("Counter %d\n", x)
		out <- x
	}
	close(out)
}

func squarer(out chan<- int, in <-chan int) {
	fmt.Println("squarer")
	for v := range in {
		fmt.Printf("Squarer %d\n", v*v)
		out <- v * v
	}
	close(out)
}

func printer(in <-chan int) {
	fmt.Println("printer")
	for v := range in {
		fmt.Printf("Printer %d\n", v)
		fmt.Println(v)
	}
}

func main() {
	naturals := make(chan int)
	squares := make(chan int)

	go counter(naturals)
	go squarer(squares, naturals)
	printer(squares)
}

//!-