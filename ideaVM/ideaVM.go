package ideaVM

// the thinky boy's thinking part
type IdeaVM struct {
	stack        []int64
	pointInStack int //golang arrays can only be 2^32 long, so might as well use the extra space to store a -1, for the end

	inputs  []int64 //the input parameters. Not changeable by the operations
	storage []int64 //a storage middle point. Changeable from within
	outputs []int64 //the output channels

	code        []operation
	pointInCode uint64
	stepsTaken  uint64
	maxSteps    uint64
}

func NewVM(inputs []int64, outputSize uint) *IdeaVM {
	return &IdeaVM{
		stack:        make([]int64, 100), //preallocating space
		pointInStack: -1,

		inputs:  inputs,
		storage: make([]int64, 100),        //preallocating space
		outputs: make([]int64, outputSize), //limit the output size
	}
}

func (vm *IdeaVM) RunModel(m *Model) []int64 {
	vm.pointInCode = 0
	vm.pointInStack = -1
	vm.code = m.operations
	return vm.Run()
}

// runs until it cant, and returns its answer array
func (vm *IdeaVM) Run() []int64 {
	for vm.canStep() {
		vm.unsafeStep()
	}
	return vm.outputs
}
func (vm *IdeaVM) Step() {
	if !vm.canStep() {
		return //skip if we cant step
	}
	vm.unsafeStep()

}
func (vm *IdeaVM) unsafeStep() {
	vm.stepsTaken++ //we can take another step, so lets log that
	vm.code[vm.pointInCode].Do(vm)
}
func (vm *IdeaVM) canStep() bool {
	if vm.maxSteps != 0 && vm.stepsTaken > vm.maxSteps {
		//stops if its passed its limit
		return false
	}
	if int(vm.pointInCode) >= len(vm.code) {
		//if it's gone past the limit for the code, stop.
		return false
	}
	return true
}

func (vm *IdeaVM) popFromStack() int64 {
	//lets return 0 if it cant return otherwise
	if vm.pointInStack == -1 {
		return 0
	}
	ans := vm.stack[vm.pointInStack]
	vm.pointInStack -= 1

	return ans
}
func (vm *IdeaVM) pushToStack(val int64) {
	vm.pointInStack += 1
	for vm.pointInStack >= len(vm.stack) {
		//the point in the stack has surpassed the size of the stack, so well expand it.
		//we do this instead of always changing the stack to try and minimize the memory reallocations needed
		vm.stack = append(vm.stack, 0)
	}
	vm.stack[vm.pointInStack] = val
}
