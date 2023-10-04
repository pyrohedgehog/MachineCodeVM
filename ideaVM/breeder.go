package ideaVM

import "math/rand"

type Breeder struct {
	//a base breeder, this is the element that would most likely be updated
	spawn []*Model
	//mutation rate is 0-1 for what percentage of the operations in a model should be mutated
	//TODO: use!
	// mutationRate float64
	evaluator EvaluatorImpl
}

// spawn new models
func (b *Breeder) SpawnModels(modelCount int, initialStepCount int) {
	//make modelCount models, each time, give it the initial step count number of operations
	models := make([]*Model, modelCount)
	for i := 0; i < modelCount; i++ {
		operations := make([]operation, initialStepCount)
		for j := 0; j < initialStepCount; j++ {
			operations[j] = GetNewOperationAtRandom()
		}
		models[i] = &Model{
			operations: operations,
		}
	}
	b.spawn = models
}

func (b *Breeder) EvaluateModels() (topEvaluation float64) {
	b.spawn = b.evaluator.OrderModels(b.spawn)
	return b.evaluator.EvaluateIndividual(b.spawn[0])
}

func (b *Breeder) CreateNextGeneration() {
	//some amount should be an amalgamation, some should be mutated pure copies of the top performer
	//lets remove the bottom half,

	//actions upon operation selection should be
	//add(before the current operation, interject a new one, selected at random)
	//subtract(remove the operation selected)
	//change(change the selected operation to a new one. If its a constant, 50/50 change it to a new value, or new constant)
	//if the len(operations)+1 was selected, we ignore it. This acts as a way to copy highly effective models in pure form
	halfSpawnSize := len(b.spawn) / 2
	for i := 0; i < halfSpawnSize; i++ {
		b.spawn[halfSpawnSize+i] = b.mutateModel(b.spawn[i])
	}
}

func (b *Breeder) mutateModel(m *Model) *Model {
	operationSelected := rand.Intn((len(m.operations) + 1))
	if len(m.operations) <= operationSelected {
		return m
	}
	mod := &Model{
		operations: m.operations,
	} //TODO: check if this is creating a new copy properly

	switch selectionType := rand.Intn(4); selectionType {
	case 0:
		//add

		selection := GetNewOperationAtRandom()
		mod.operations = append(mod.operations[:operationSelected], selection)
		mod.operations = append(mod.operations, mod.operations[operationSelected:]...)
	case 1:
		//subtract
		mod.operations = append(mod.operations[:operationSelected], mod.operations[:operationSelected+1]...)
	case 2:
		//change
		if mod.operations[operationSelected].GetType() == Const &&
			rand.Float64() > 0.5 {
			//TODO: change that to a percentage chance of changing the value
			op := mod.operations[operationSelected].(opConst)
			op.value += rand.Int63n(1<<32-1) * (int64(rand.Intn(2) * -1))
			mod.operations[operationSelected] = op
		} else {
			mod.operations[operationSelected] = GetNewOperationAtRandom()
		}

	}

	return mod
}
