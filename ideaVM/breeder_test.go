package ideaVM

import (
	"fmt"
	"math"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBreederSpawnGeneration(t *testing.T) {
	doFancyOutputs := false
	breeder := &Breeder{
		evaluator: NewSimpleEvaluator(addTwoGenerateTests(50)),
	}
	breeder.SpawnModels(1000, 100)
	//generating 100 spawn, all 9 of the operations should showup at least once.
	operationTypes := map[OpCode]int{
		Add:         0,
		Mul:         0,
		LSInput:     0,
		RSInput:     0,
		GetInput:    0,
		LSOutput:    0,
		RSOutput:    0,
		WriteOutput: 0,
		Const:       0,
	}
	totalOperations := 0
	for _, v := range breeder.spawn {
		for _, op := range v.operations {
			operationTypes[op.GetType()] += 1
			totalOperations += 1
		}
	}
	if doFancyOutputs {
		//we now have a recording of how often each one shows up.
		fmt.Println("frequency of each sign graphed below")
		fmt.Printf("Add      : %v\n", operationTypes[Add])
		fmt.Printf("Mul      : %v\n", operationTypes[Mul])
		fmt.Printf("RSInput  : %v\n", operationTypes[LSInput])
		fmt.Printf("LSInput  : %v\n", operationTypes[RSInput])
		fmt.Printf("GetInput : %v\n", operationTypes[GetInput])
		fmt.Printf("LSOut    : %v\n", operationTypes[LSOutput])
		fmt.Printf("RSOut    : %v\n", operationTypes[RSOutput])
		fmt.Printf("WriteOut : %v\n", operationTypes[WriteOutput])
		fmt.Printf("Const    : %v\n", operationTypes[Const])
	}
	operationFloats := map[OpCode]float64{
		Add:         0,
		Mul:         0,
		LSInput:     0,
		RSInput:     0,
		GetInput:    0,
		LSOutput:    0,
		RSOutput:    0,
		WriteOutput: 0,
		Const:       0,
	}
	maxPercent := float64(0)
	minPercent := float64(1000)
	for op, val := range operationTypes {
		operationFloats[op] = float64(val) / float64(totalOperations) * 100
		if operationFloats[op] > maxPercent {
			maxPercent = operationFloats[op]
		}
		if operationFloats[op] < minPercent {
			minPercent = operationFloats[op]
		}
	}
	percentVariance := maxPercent - minPercent
	assert.Less(t,
		percentVariance, float64(5),
		"distribution appears to be too far off",
	)
}
func TestBreedingForAddTwo(t *testing.T) {
	breeder := runTestOnSimpleEvaluator(t, addTwoGenerateTests(100), 100, 5)
	topScore := breeder.evaluator.EvaluateIndividual(breeder.spawn[0])
	fmt.Printf("A top performer has been evaluated at %v%% accuracy\n", topScore*100)
	fmt.Println("congratulations!!!")
}
func TestBreedingForReturnNumber(t *testing.T) {
	breeder := runTestOnSimpleEvaluator(t, returnNumberGenerateTests(100), 100, 3)
	topScore := breeder.evaluator.EvaluateIndividual(breeder.spawn[0])
	fmt.Printf("A top performer has been evaluated at %v%% accuracy\n", topScore*100)
	fmt.Printf("congratulations!!!")
}

func runTestOnSimpleEvaluator(t *testing.T, testSet []SimpleTestSet, generationSize int, modelOperationSize int) *Breeder {
	eval := NewSimpleEvaluator(testSet)
	breeder := &Breeder{
		evaluator: eval,
	}
	genNumber := 0
	breeder.SpawnModels(generationSize, modelOperationSize)
	//DELETE AFTER EVALUATING IF THIS WORKS
	// breeder.spawn[0] = &Model{operations: []operation{
	// 	opGetInput{},
	// 	opRSInputs{},
	// 	opGetInput{},
	// 	opAdd{},
	// 	opWriteOutput{},
	// }}
	//DELETE!!!
	topScore := breeder.EvaluateModels()
	lastScore := topScore
	fmt.Printf("Evaluated Spawn to have a top performer at %v%% accuracy.\n", topScore*100)
	for topScore <= 0.99 {
		genLog := math.Log10(float64(genNumber))
		if genNumber < 10 ||
			genNumber%int(math.Pow10(int(genLog))) == 0 ||
			lastScore < topScore {
			fmt.Printf("gen:%v\t\ttop:%v%%\n", genNumber, topScore*100)
		}
		breeder.CreateNextGeneration()
		lastScore = topScore
		topScore = breeder.EvaluateModels()
		genNumber += 1
	}
	return breeder
}

func addTwoGenerateTests(testCount int) []SimpleTestSet {
	set := make([]SimpleTestSet, testCount)
	for i := 0; i < testCount; i++ {
		ins := []int64{rand.Int63(), rand.Int63()}
		set[i].inputs = ins
		set[i].expectedResults = []int64{ins[0] + ins[1]}
	}

	return set
}

func returnNumberGenerateTests(testCount int) []SimpleTestSet {
	set := make([]SimpleTestSet, testCount)
	for i := 0; i < testCount; i++ {
		ins := []int64{rand.Int63()}
		set[i].inputs = ins
		set[i].expectedResults = []int64{ins[0]}
	}

	return set
}
