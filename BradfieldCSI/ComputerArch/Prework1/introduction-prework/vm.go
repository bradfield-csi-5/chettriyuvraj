package vm

import (
	"log"
)

const (
	Load  = 0x01
	Store = 0x02
	Add   = 0x03
	Sub   = 0x04
	Halt  = 0xff
)

// Stretch goals
const (
	Addi = 0x05
	Subi = 0x06
	Jump = 0x07
	Beqz = 0x08
)

const (
	pc = iota
	r1
	r2
)

// Given a 256 byte array of "memory", run the stored program
// to completion, modifying the data in place to reflect the result
//
// The memory format is:
//
// 00 01 02 03 04 05 06 07 08 09 0a 0b 0c 0d 0e 0f ... ff
// __ __ __ __ __ __ __ __ __ __ __ __ __ __ __ __ ... __
// ^==DATA===============^ ^==INSTRUCTIONS==============^

func compute(memory []byte) {

	registers := [3]byte{8, 0, 0} // PC, R1 and R2

	// Keep looping, like a physical computer's clock
	for {

		/* Fetch */
		opcode := memory[registers[pc]]

		/* Decode */
		var operands uint16
		var maxOffset byte

		/* Using two switches because doing this in a single switch with fallthrough is restricted in Go */
		switch opcode {
		case Halt:
			maxOffset = 0x00
		case Jump:
			maxOffset = 0x01
			operands = decodeOperands(memory, registers, maxOffset)
		default:
			maxOffset = 0x02
			operands = decodeOperands(memory, registers, maxOffset)
		}

		/* Execute */
		switch opcode {
		case Load:
			loadProcedure(operands, memory, &registers)
		case Store:
			storeProcedure(operands, memory, &registers)
		case Add:
			addProcedure(operands, &registers)
		case Sub:
			subProcedure(operands, &registers)
		case Addi:
			addiProcedure(operands, &registers)
		case Subi:
			subiProcedure(operands, &registers)
		case Jump:
			jumpProcedure(operands, &registers)
		case Beqz:
			beqzProcedure(operands, &registers)
		case Halt:
			break
		default:
			log.Fatalf("invalid opcode %x", opcode)
			continue
		}

		/* Increase program counter */
		if opcode != Jump {
			registers[pc] += maxOffset + 1
		}

		if opcode == Halt {
			break
		}
	}

}

func decodeOperands(memory []byte, registers [3]byte, maxOffset byte) (operands uint16) {
	for offset := maxOffset; offset > 0x00; offset-- {
		operands |= (uint16(memory[registers[pc]+offset]) << (8 * (maxOffset - offset))) /* A little convoluted - we are basically grabbing the operands from right to left in a 16 bit value */
	}
	return operands
}

/* Load value at given address into given register */
func loadProcedure(operands uint16, memory []byte, registers *[3]byte) {
	address := byte(0x00 | operands)
	operands >>= 8
	register := byte(0x00 | operands)
	registers[register] = memory[address]
}

/* Store value at given register into given address */
func storeProcedure(operands uint16, memory []byte, registers *[3]byte) {
	a1 := byte(0x00 | operands)
	operands >>= 8
	r1 := byte(0x00 | operands)
	memory[a1] = registers[r1]
}

/* Add two register values and store into destination register*/
func addProcedure(operands uint16, registers *[3]byte) {
	r2 := byte(0x00 | operands)
	operands >>= 8
	r1 := byte(0x00 | operands) /* Destination register */
	registers[r1] = registers[r1] + registers[r2]
}

/* Add two register values and store into destination register*/
func subProcedure(operands uint16, registers *[3]byte) {
	r2 := byte(0x00 | operands)
	operands >>= 8
	r1 := byte(0x00 | operands)
	registers[r1] = registers[r1] - registers[r2]
}

/* Add one op and one register value and store into register*/
func addiProcedure(operands uint16, registers *[3]byte) {
	op1 := byte(0x00 | operands)
	operands >>= 8
	r1 := byte(0x00 | operands)
	registers[r1] = registers[r1] + op1
}

/* Add one op from one register value and store into register*/
func subiProcedure(operands uint16, registers *[3]byte) {
	op1 := byte(0x00 | operands)
	operands >>= 8
	r1 := byte(0x00 | operands)
	registers[r1] = registers[r1] - op1
}

/* Jump program counter to desired location in memory - can use only uint8 operand here*/
func jumpProcedure(operands uint16, registers *[3]byte) {
	jumpAddr := byte(0x00 | operands)
	registers[pc] = jumpAddr
}

/* Jump program counter to desired location in memory*/
func beqzProcedure(operands uint16, registers *[3]byte) {
	offset := byte(0x00 | operands)
	operands >>= 8
	r1 := byte(0x00 | operands)
	if registers[r1] == 0x00 {
		registers[pc] += offset
	}
}
