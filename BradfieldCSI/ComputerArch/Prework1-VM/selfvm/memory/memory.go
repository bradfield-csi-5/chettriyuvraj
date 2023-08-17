/* Package to simulate main memory in a computer */

package memory

import "fmt"

const (
	MEMORYSIZE        = 1 << 8
	DATABOUNDARY byte = 0x07
)

/* Memory has a (logical) split between different regions for data, output(first byte of data) and instructions */
type Memory [MEMORYSIZE]byte

func NewMemory() *Memory {
	memory := Memory([MEMORYSIZE]byte{})
	return &memory
}

/* Start and end signify the boundary to write within */
func (memory *Memory) writeToMemory(b []byte, start byte, end byte) (int, error) {
	if len(b) > int(end-start+1) {
		return 0, fmt.Errorf("not enough memory to write to")
	}

	for i, j := start, 0; j < len(b); i, j = i+1, j+1 {
		(*memory)[i] = b[j]
	}

	return len(b), nil
}

func (memory *Memory) WriteInstructions(b []byte, start byte) (int, error) {
	return memory.writeToMemory(b, start, MEMORYSIZE-1)
}

func (memory *Memory) WriteData(b []byte, start byte) (int, error) {
	return memory.writeToMemory(b, start, DATABOUNDARY)
}
