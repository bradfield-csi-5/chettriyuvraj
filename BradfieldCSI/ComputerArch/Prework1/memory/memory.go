/* Package to simulate main memory in a computer */

package memory

import "fmt"

type Memory []byte

func (m *Memory) Write(b []byte) (int, error) {
	if len(*m) < len(b) {
		return 0, fmt.Errorf("not enough memory to write to")
	}

	for i, j := 0, 0; j < len(b); i, j = i+1, j+1 {
		(*m)[i] = b[j]
	}
	return len(b), nil
}

type MainMemory struct {
	Data         Memory
	Instructions Memory
	Output       Memory /* First x bytes of data slice will be assigned to output */
}

func NewMainMemory(instructionMemorySize int, dataMemorySize int, outputMemorySize int) (*MainMemory, error) {
	if outputMemorySize >= dataMemorySize { /* OutputMemorySize must be atleast 1 byte lesser than dataMemory */
		return nil, fmt.Errorf("output memory greater than or equal to data memory")
	}

	mainMemorySize := instructionMemorySize + dataMemorySize
	allocatedMemory := make([]byte, mainMemorySize)

	mainMemory := MainMemory{Output: allocatedMemory[:outputMemorySize], Data: allocatedMemory[:dataMemorySize], Instructions: allocatedMemory[dataMemorySize:]}
	return &mainMemory, nil
}
