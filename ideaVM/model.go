package ideaVM

// a model stores the actual operation steps it uses.
type Model struct {
	operations []operation
}
type valuedModel struct {
	score float64
	model *Model
}

func GetModelFromStorable(vals []int64) *Model {
	ops := []operation{}
	for i := 0; i < len(vals); {
		switch vals[i] {
		case int64(Add):
			ops = append(ops, opAdd{})
			i += 1
		case int64(Mul):
			ops = append(ops, opMul{})
			i += 1
		case int64(GetInput):
			ops = append(ops, opGetInput{})
			i += 1
		case int64(WriteOutput):
			ops = append(ops, opWriteOutput{})
			i += 1
		case int64(Const):
			//we need to get the value, so this needs a custom handling
			ops = append(ops, opConst{
				value: vals[i+1],
			})
			i += 2

		}
	}

	return &Model{
		operations: ops,
	}
}
func (m *Model) GetStorable() []int64 {
	//take all of the operations, if they have a value, store that after the type
	ansBytes := []int64{}
	for _, op := range m.operations {
		ansBytes = append(ansBytes, int64(op.GetType()))
		switch op.GetType() {
		//only handle the edge cases that have values here
		case Const:
			ansBytes = append(ansBytes, op.(opConst).value)

		}
	}
	return ansBytes
}
