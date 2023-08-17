// package to simulate the cpu of a computer

/* Plan
1. Control Unit
2. ALU
3. Registers
*/

package cpu

import (
	"fmt"
	"virtualcomputer/memory"
)

const (
	pc byte = iota
	r1
	r2
)

type CPU struct {
	controlUnit ControlUnit
	registers   *[3]byte
	mainMemory  *memory.Memory
}

func NewCPU() CPU {
	return CPU{controlUnit: NewControlUnit(), registers: &[3]byte{memory.DATABOUNDARY + 0x01}, mainMemory: memory.NewMemory()}
}

func (cpu CPU) Fetch() (instruction uint32) {
	return cpu.controlUnit.Fetch(cpu.mainMemory, cpu.registers)
}

func (cpu CPU) DecodeAndExecute(instruction uint32) (opcode byte, err error) {
	return cpu.controlUnit.DecodeAndExecute(instruction, cpu.mainMemory, cpu.registers)
}

func (cpu CPU) WriteInstructions(b []byte, start byte) (int, error) {
	return cpu.mainMemory.WriteInstructions(b, start)
}

func (cpu CPU) WriteData(b []byte, start byte) (int, error) {
	return cpu.mainMemory.WriteData(b, start)
}

func (cpu CPU) String() string {
	return fmt.Sprintf("registers %x \n mainMemory  %x\n registers %b \n mainMemory  %b\n", cpu.registers, *(cpu.mainMemory), cpu.registers, *(cpu.mainMemory))
}

/* Load value at given address into given register */
func loadProcedure(instruction uint32, memory *memory.Memory, registers *[3]byte) {
	address := byte(0x00 | instruction)
	instruction >>= 8
	register := byte(0x00 | instruction)
	(*registers)[register] = (*memory)[address]
}

/* Store value at given register into given address */
func storeProcedure(instruction uint32, memory *memory.Memory, registers *[3]byte) {
	address := byte(0x00 | instruction)
	instruction >>= 8
	register := byte(0x00 | instruction)
	(*memory)[address] = (*registers)[register]
}

/* Add two values and store into destination register*/
func addProcedure(instruction uint32, registers *[3]byte) {
	r2 := byte(0x00 | instruction)
	instruction >>= 8
	r1 := byte(0x00 | instruction) /* Destination register */
	(*registers)[r1] = (*registers)[r1] + (*registers)[r2]
}

/* Add two values and store into destination register*/
func subProcedure(instruction uint32, registers *[3]byte) {
	r2 := byte(0x00 | instruction)
	instruction >>= 8
	r1 := byte(0x00 | instruction) /* Destination register */
	(*registers)[r1] = (*registers)[r1] - (*registers)[r2]
}
