package main

import "testing"

func BenchmarkEcho(b *testing.B) {
	for i := 0; i < b.N; i++ {
		echo()
	}
}

func BenchmarkEchoOptimized(b *testing.B) {
	for i := 0; i < b.N; i++ {
		echoOptimized()
	}
}

/*
BenchmarkEcho
PASS
ok      Chapter1        4.169s

BenchmarkEchoOptimized
PASS
ok      Chapter1        2.553s
*/
