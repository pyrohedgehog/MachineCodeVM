package ideaVM

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
)

func TestBreedingToSolution(t *testing.T) {
	eval := NewSimpleEvaluator(addTwoGenerateTests(500))
	breeder := &Breeder{
		evaluator: eval,
	}
	genNumber := 0
	breeder.SpawnModels(500, 6)
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
	fmt.Printf("A top performer has been evaluated at %v%% accuracy\n", topScore)
	fmt.Printf("congratulations!!!")
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
