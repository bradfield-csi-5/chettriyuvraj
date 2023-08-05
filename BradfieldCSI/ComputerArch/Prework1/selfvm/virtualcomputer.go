package main

import (
	"log"
	"virtualcomputer/cpu"
)

func main() {
	myCpu := cpu.NewCPU()

	/* Write to memory before starting cycle  */

	for {
		instruction := myCpu.Fetch()

		opcode, err := myCpu.DecodeAndExecute(instruction)
		if err != nil {
			log.Fatalf("error in decode-execute phase %w", err)
		}
		if opcode == 0xff {
			break
		}
	}

}
