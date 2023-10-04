package ideaVM

type MemoryReal struct {
	//memory real acts like a magnetic strip, of a fixed length that circles, but can be as large as 2^^63
	memory      []int64
	accessPoint int
}

func NewMemoryReal(memorySize int) *MemoryReal {
	mr := &MemoryReal{
		memory:      make([]int64, memorySize),
		accessPoint: 0,
	}
	return mr

}
func NewMemoryRealWith(values []int64) *MemoryReal {
	mr := MemoryReal{
		memory:      make([]int64, len(values)),
		accessPoint: 0,
	}
	copy(mr.memory, values)
	return &mr
}
func (mr *MemoryReal) GetValue() int64 {
	return mr.memory[mr.accessPoint]
}
func (mr *MemoryReal) GetMemory() []int64 {
	return mr.memory
}
func (mr *MemoryReal) WriteValue(val int64) {
	mr.memory[mr.accessPoint] = val
}
func (mr *MemoryReal) LeftShift() {
	mr.accessPoint -= 1
	if mr.accessPoint < 0 {
		mr.accessPoint = len(mr.memory) - 1
	}
}
func (mr *MemoryReal) RightShift() {
	mr.accessPoint = (mr.accessPoint + 1) % len(mr.memory)
}
