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
	vms      []*ThreadedIdeaVM
	testData []SimpleTestSet
}

func NewSimpleEvaluator(testData []SimpleTestSet) *SimpleEvaluator {
	e := SimpleEvaluator{
		vms:      make([]*ThreadedIdeaVM, len(testData)),
		testData: testData,
	}
	for i, test := range testData {
		v := NewThreadedVM(test.inputs, uint(len(test.expectedResults)))
		v.vm.maxSteps = 1500 //TODO: make this a param
		e.vms[i] = v
	}

	return &e
}

type SimpleTestSet struct {
	inputs          []int64
	expectedResults []int64
}

// var maxInt float64 = (1 << 64) - 1

func (e *SimpleEvaluator) OrderModels(models []*Model) []*Model {
	//see how close it is to giving the target answer. should be run more than once per evaluation with multiple inputs
	type valuedModel struct {
		score float64
		model *Model
	}
	valuedModelsChan := make(chan valuedModel)
	//TODO: i started having a migraine while writing this, i need to clean it.
	for _, v := range models {
		go func(v *Model) {
			mod := valuedModel{
				score: e.EvaluateIndividual(v),
				model: v,
			}
			valuedModelsChan <- mod
		}(v)
	}
	//get that out from the channel so we can sort it
	valuedModels := make([]valuedModel, len(models))
	for i := 0; i < len(models); i++ {
		valuedModels[i] = <-valuedModelsChan
	}
	sort.Slice(valuedModels, func(i, j int) bool {
		return valuedModels[i].score > valuedModels[j].score
	})

	//then get that out from the valued models version into a normal array of models
	orderedAns := make([]*Model, len(models))
	for x, v := range valuedModels {
		orderedAns[x] = v.model
	}
	return orderedAns
}

func (e *SimpleEvaluator) EvaluateIndividual(m *Model) float64 {
	scoreChan := make(chan float64)
	//TODO: think about how we evaluate a lot more here, even in the simplest version
	for j, vm := range e.vms {
		go func() {
			var roundScore float64 = 0
			result := vm.RunModel(m)
			expected := e.testData[j].expectedResults
			//for each of the results, find how close they are to the correct answer, and then evaluate it
			for resultPoint := 0; resultPoint < len(result); resultPoint++ {
				val := expected[resultPoint] - result[resultPoint]
				// score += float64(1) - (math.Abs(float64(val)) / maxInt)
				if val == 0 {
					roundScore += 1
				} else {
					roundScore += float64(1) / math.Log10(math.Abs(float64(val)))
				}
			}
			scoreChan <- roundScore / float64(len(expected))
		}()
	}
	score := float64(0)
	for i := 0; i < len(e.vms); i++ {
		foo := <-scoreChan
		score += foo
	}
	score = score / float64(len(e.testData))
	return score
}
