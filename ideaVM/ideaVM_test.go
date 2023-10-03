package ideaVM

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVM(t *testing.T) {
	vm := NewVM([]int64{1, 2}, 2)
	vm.pushToStack(1)
	assert.Equal(t,
		0, vm.pointInStack,
		"wrong point in stack",
	)
	assert.Equal(t,
		int64(1), vm.popFromStack(),
		"stack returned wrong val",
	)
	assert.Equal(t,
		-1, vm.pointInStack,
		"wrong point in stack",
	)
	addOp := opAdd{}
	vm.pushToStack(1)
	vm.pushToStack(2)
	addOp.Do(vm)

	assert.Equal(t,
		int64(3), vm.popFromStack(),
		"stack returned wrong val",
	)
}

func TestIdealAddFunction(t *testing.T) {
	testParameters := [][]int64{{1, 2}, {3, 4}, {-1, 5}}
	vm := NewVM([]int64{1, 2}, 1)
	vm.code = []operation{
		opGetInput{},
		opConst{value: 1},
		opGetInput{},
		opAdd{},
		opConst{value: 0},
		opWriteOutput{},
	}
	for _, param := range testParameters {
		vm.inputs = param
		vm.pointInCode = 0
		vm.pointInStack = -1
		ans := vm.Run()
		assert.Equal(t,
			[]int64{param[0] + param[1]}, ans,
			"could not calculate answer correctly",
		)
	}

}
