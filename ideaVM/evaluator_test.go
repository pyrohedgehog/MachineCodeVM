package ideaVM

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleEvaluations(t *testing.T) {
	//test based on a function that adds two numbers together
	addTwo := Model{operations: []operation{
		opGetInput{},
		opRSInputs{},
		opGetInput{},
		opAdd{},
		opWriteOutput{},
	}}
	multiplyTwo := Model{operations: []operation{
		opGetInput{},
		opRSInputs{},
		opGetInput{},
		opMul{},
		opWriteOutput{},
	}}

	addTwoTestData := []SimpleTestSet{
		{inputs: []int64{100, 200},
			expectedResults: []int64{300}},
		{inputs: []int64{300, 400},
			expectedResults: []int64{700}},
	}
	eval := NewSimpleEvaluator(addTwoTestData)
	addTwoEval := eval.EvaluateIndividual(&addTwo)
	mulTwoEval := eval.EvaluateIndividual(&multiplyTwo)
	assert.Greater(t,
		addTwoEval, mulTwoEval,
		"evaluator evaluated the wrong answer as more correct",
	)
	foo := []*Model{&addTwo, &multiplyTwo}
	bar := eval.OrderModels(foo)
	assert.Equal(t,
		eval.EvaluateIndividual(bar[0]), float64(1),
		"not ordering items properly",
	)

}
