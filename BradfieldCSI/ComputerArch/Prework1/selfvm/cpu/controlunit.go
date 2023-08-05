package cpu

import (
	"fmt"
	"virtualcomputer/memory"
)

type ControlUnit struct{}

const (
	_ byte = iota
	load
	store
	add
	sub
	halt = 0xff
)

func NewControlUnit() ControlUnit {
	return ControlUnit{}
}

func (controlUnit ControlUnit) Fetch(memory *memory.Memory, registers *[3]byte) (instruction uint32) {
	for offset := byte(0x00); offset < 3; offset, instruction = offset+1, instruction<<8 {
		fmt.Printf("or %x\n", uint32((*memory)[(*registers)[pc]+offset]))
		instruction |= uint32((*memory)[(*registers)[pc]+offset])
		fmt.Printf("i %x\n", instruction)

	}
	instruction >>= 8
	fmt.Printf("Instruction %x\n", instruction)
	return instruction
}

func (controlUnit ControlUnit) DecodeAndExecute(instruction uint32, memory *memory.Memory, registers *[3]byte) (opcode byte, err error) {
	opcode = byte(0x00 | instruction>>16)

	fmt.Printf("opcode %x\n", opcode)

	switch opcode {
	case load:
		loadProcedure(instruction, memory, registers)
	case store:
		storeProcedure(instruction, memory, registers)
	case add:
		addProcedure(instruction, registers)
	case sub:
		subProcedure(instruction, registers)
	case halt:
		break
	default:
		return opcode, fmt.Errorf("invalid opcode")
	}

	if opcode == halt {
		(*registers)[pc] += 0x01
	} else {
		(*registers)[pc] += 0x03
	}

	return opcode, nil
}
