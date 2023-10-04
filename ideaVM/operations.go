package ideaVM

import "math/rand"

type operation interface {
	//the methods an opcode can do. Any opcode
	Do(*IdeaVM)
	GetType() OpCode
}

var possibilities = []operation{
	opAdd{}, opMul{},
	opLSInputs{}, opRSInputs{}, opGetInput{},
	opLSOutputs{}, opRSOutputs{}, opWriteOutput{},
	opConst{},
} //all the possible options it can generate with
func GetNewOperationAtRandom() operation {
	op := possibilities[rand.Intn(len(possibilities))]
	var val int64
	if op.GetType() == Const {
		switch rand.Intn(4) {
		case 0:
			val = 0
		case 1:
			val = 1
		case 2: //negative
			val = rand.Int63()
			val *= -1
		case 3:
			val = rand.Int63()
		}
		op = opConst{
			value: val,
		}
	}
	return op
}

type opAdd struct{}

func (op opAdd) Do(vm *IdeaVM) {
	vm.pushToStack(
		vm.popFromStack() + vm.popFromStack(),
	)
	vm.pointInCode++
}
func (op opAdd) GetType() OpCode {
	return Add
}

type opMul struct{}

func (op opMul) Do(vm *IdeaVM) {
	vm.pushToStack(
		vm.popFromStack() * vm.popFromStack(),
	)
	vm.pointInCode++
}
func (op opMul) GetType() OpCode { return Mul }

type opLSInputs struct{}

func (op opLSInputs) Do(vm *IdeaVM) {
	vm.inputs.LeftShift()
	vm.pointInCode++
}
func (op opLSInputs) GetType() OpCode { return LSInput }

type opRSInputs struct{}

func (op opRSInputs) Do(vm *IdeaVM) {
	vm.inputs.RightShift()
	vm.pointInCode++
}
func (op opRSInputs) GetType() OpCode { return RSInput }

type opGetInput struct{}

func (op opGetInput) Do(vm *IdeaVM) {
	vm.pushToStack(vm.inputs.GetValue())
	vm.pointInCode++
}
func (op opGetInput) GetType() OpCode { return GetInput }

type opLSOutputs struct{}

func (op opLSOutputs) Do(vm *IdeaVM) {
	vm.outputs.LeftShift()
	vm.pointInCode++
}
func (op opLSOutputs) GetType() OpCode { return LSOutput }

type opRSOutputs struct{}

func (op opRSOutputs) Do(vm *IdeaVM) {
	vm.outputs.RightShift()
	vm.pointInCode++
}
func (op opRSOutputs) GetType() OpCode { return RSOutput }

type opWriteOutput struct{}

func (op opWriteOutput) Do(vm *IdeaVM) {
	writeVal := vm.popFromStack()
	vm.outputs.WriteValue(writeVal)
	vm.pointInCode++
}
func (op opWriteOutput) GetType() OpCode { return WriteOutput }

type opConst struct {
	value int64
}

func (op opConst) Do(vm *IdeaVM) {
	vm.pushToStack(op.value)
	vm.pointInCode++
}
func (op opConst) GetType() OpCode { return Const }
