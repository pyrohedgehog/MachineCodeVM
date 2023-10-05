package ideaVM

import (
	"fmt"
	"math/rand"
)

type Breeder struct {
	//a base breeder, this is the element that would most likely be updated
	spawn              []*Model
	initialModelLength int
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
	b.initialModelLength = initialStepCount
}

func (b *Breeder) EvaluateModels() (topEvaluation float64) {
	b.spawn = b.evaluator.OrderModels(b.spawn)
	return b.evaluator.EvaluateIndividual(b.spawn[0])
}
func (b *Breeder) EvaluateAndBreed(skipIfAchieved float64) (topGrade float64, bottomGrade float64, average float64) {
	min, max, avg, valued := b.evaluator.GetEvaluationProperties(b.spawn)
	topGrade = max
	bottomGrade = min
	average = avg
	if max >= skipIfAchieved {
		//so you can run this in a loop until it reaches this value
		return
	}
	children := []*Model{}
	if max-min <= float64(0.0001) {
		//there's a lot of homogony in here, so lets just cutoff a large part, and throw in some new genes.
		valued = valued[len(valued)/4:]
		for i := 0; i < len(b.spawn)/5; i++ {
			//20% of the children will be entirely new.
			freshSpawn := &Model{
				operations: make([]operation, b.initialModelLength),
			}
			for j := 0; j < b.initialModelLength; j++ {
				freshSpawn.operations[j] = GetNewOperationAtRandom()
			}
			children = append(children, freshSpawn)
		}
		avg -= 0.0001
	}
	for _, val := range valued {
		if val.score >= avg && len(children) <= len(b.spawn) {
			children = append(children, val.model)
		}
	}
	for len(children) < len(b.spawn) {
		parent := children[rand.Intn(len(children))]
		children = append(children, b.mutateModel(parent))
	}

	b.spawn = children
	return
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
	operations := m.operations
	//TODO: check if this is creating a new copy properly

	switch selectionType := rand.Intn(4); selectionType {
	case 0:
		//add

		selection := GetNewOperationAtRandom()
		operations = append(operations[:operationSelected], selection)
		operations = append(operations, operations[operationSelected:]...)
	case 1:
		//subtract
		operations = append(operations[:operationSelected], operations[:operationSelected+1]...)
	case 2:
		//change
		if operations[operationSelected].GetType() == Const &&
			rand.Float64() > 0.5 {
			//TODO: change that to a percentage chance of changing the value
			op := operations[operationSelected].(opConst)
			op.value += rand.Int63n(1<<32-1) * (int64(rand.Intn(2) * -1))
			operations[operationSelected] = op
		} else {
			operations[operationSelected] = GetNewOperationAtRandom()
		}

	}

	return &Model{
		operations: operations,
	}
}
