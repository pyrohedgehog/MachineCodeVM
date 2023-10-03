package ideaVM

import (
	"math"
	"sort"
)

type EvaluatorImpl interface {
	//take an array of models, evaluate them, and return the models in order of decreasing accuracy
	OrderModels([]*Model) []*Model
	EvaluateIndividual(*Model) float64
}
type SimpleEvaluator struct {
	vms      []*IdeaVM
	testData []SimpleTestSet
}

func NewSimpleEvaluator(testData []SimpleTestSet) *SimpleEvaluator {
	e := SimpleEvaluator{
		vms:      make([]*IdeaVM, len(testData)),
		testData: testData,
	}
	for i, test := range testData {
		v := NewVM(test.inputs, uint(len(test.expectedResults)))
		v.maxSteps = 1500 //TODO: make this a param
		e.vms[i] = v
	}

	return &e
}

type SimpleTestSet struct {
	inputs          []int64
	expectedResults []int64
}

var maxInt float64 = (1 << 64) - 1

func (e *SimpleEvaluator) OrderModels(models []*Model) []*Model {
	//see how close it is to giving the target answer. should be run more than once per evaluation with multiple inputs
	type valuedModel struct {
		score float64
		model *Model
	}
	valuedModels := make([]valuedModel, len(models))
	//TODO: i started having a migraine while writing this, i need to clean it.
	for i, v := range models {
		mod := valuedModel{
			score: e.EvaluateIndividual(v),
			model: v,
		}
		valuedModels[i] = mod
	}
	sort.Slice(valuedModels, func(i, j int) bool {
		return valuedModels[i].score > valuedModels[j].score
	})
	orderedAns := make([]*Model, len(models))
	for x, v := range valuedModels {
		orderedAns[x] = v.model
	}
	return orderedAns
}

func (e *SimpleEvaluator) EvaluateIndividual(m *Model) float64 {
	score := float64(0)
	//TODO: think about how we evaluate a lot more here, even in the simplest version
	for j, vm := range e.vms {
		result := vm.RunModel(m)
		expected := e.testData[j].expectedResults
		//for each of the results, find how close they are to the correct answer, and then evaluate it
		for resultPoint := 0; resultPoint < len(result); resultPoint++ {
			val := expected[resultPoint] - result[resultPoint]
			score += float64(1) - (math.Abs(float64(val)) / maxInt)
		}
		score = score / float64(len(expected))
	}
	score = score / float64(len(e.testData))
	return score
}
