package ideaVM

type operation interface {
	//the methods an opcode can do. Any opcode
	Do(*IdeaVM)
	GetType() OpCode
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
func (op opMul) GetType() OpCode {
	return Mul
}

type opGetInput struct{}

func (op opGetInput) Do(vm *IdeaVM) {
	inputPoint := vm.popFromStack()
	if len(vm.inputs) <= int(inputPoint) || inputPoint < 0 {
		vm.pushToStack(0)
	} else {
		vm.pushToStack(vm.inputs[inputPoint])
	}
	vm.pointInCode++
}
func (op opGetInput) GetType() OpCode {
	return GetInput
}

type opWriteOutput struct{}

func (op opWriteOutput) Do(vm *IdeaVM) {
	writePoint := vm.popFromStack() % int64(len(vm.outputs))
	writeVal := vm.popFromStack()
	vm.outputs[writePoint] = writeVal
	vm.pointInCode++
}
func (op opWriteOutput) GetType() OpCode {
	return WriteOutput
}

type opConst struct {
	value int64
}

func (op opConst) Do(vm *IdeaVM) {
	vm.pushToStack(op.value)
	vm.pointInCode++
}
func (op opConst) GetType() OpCode {
	return Const
}
