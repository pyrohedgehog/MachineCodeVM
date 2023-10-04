package ideaVM

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
)

func TestBreedingForAddTwo(t *testing.T) {
	breeder := runTestOnSimpleEvaluator(t, addTwoGenerateTests(100), 10000, 8)
	topScore := breeder.evaluator.EvaluateIndividual(breeder.spawn[0])
	fmt.Printf("A top performer has been evaluated at %v%% accuracy\n", topScore)
	fmt.Printf("congratulations!!!")
}
func TestBreedingForReturnNumber(t *testing.T) {
	breeder := runTestOnSimpleEvaluator(t, returnNumberGenerateTests(100), 10000, 3)
	topScore := breeder.evaluator.EvaluateIndividual(breeder.spawn[0])
	fmt.Printf("A top performer has been evaluated at %v%% accuracy\n", topScore)
	fmt.Printf("congratulations!!!")
}

func runTestOnSimpleEvaluator(t *testing.T, testSet []SimpleTestSet, generationSize int, modelOperationSize int) *Breeder {
	eval := NewSimpleEvaluator(testSet)
	breeder := &Breeder{
		evaluator: eval,
	}
	genNumber := 0
	breeder.SpawnModels(generationSize, modelOperationSize)
	topScore := breeder.EvaluateModels()
	lastScore := topScore
	fmt.Printf("Evaluated Spawn to have a top performer at %v%% accuracy.\n", topScore)
	for topScore <= 0.99 {
		genLog := math.Log10(float64(genNumber))
		if genNumber < 10 ||
			genNumber%int(math.Pow10(int(genLog))) == 0 ||
			lastScore != topScore {
			fmt.Printf("gen:%v\t\ttop:%v%%\n", genNumber, topScore)
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
