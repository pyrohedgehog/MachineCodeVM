package ideaVM

import "sync"

type ThreadedIdeaVM struct {
	vm   IdeaVM
	lock sync.Mutex
}

func NewThreadedVM(inputs []int64, outputSize uint) *ThreadedIdeaVM {
	tvm := ThreadedIdeaVM{
		vm: *NewVM(inputs, outputSize),
	}

	return &tvm
}
func (tvm *ThreadedIdeaVM) RunModel(m *Model) []int64 {
	tvm.lock.Lock()
	ans := tvm.vm.RunModel(m)
	tvm.lock.Unlock()
	return ans
}
