package ideaVM

type OpCode int64

const (
	Add OpCode = iota
	Mul
	GetInput
	WriteOutput
	Const
)
