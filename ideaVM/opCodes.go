package ideaVM

type OpCode int64

const (
	Add OpCode = iota
	Mul
	LSInput
	RSInput
	GetInput
	LSOutput
	RSOutput
	WriteOutput
	Const
)
