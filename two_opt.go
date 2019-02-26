package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Node struct {
	NodeNumber int
}

type TSP struct {
	Tour           []Node
	DistanceMatrix [][]int
	NumberOfNodes  int
}

func (t *TSP) calculate_objective() int {
	// initialize cost of going from last index to first index
	node1 := t.Tour[0].NodeNumber
	node2 := t.Tour[len(t.Tour)-1].NodeNumber
	cost := t.DistanceMatrix[node1][node2]
	// iterate through on a slice to calculate costs
	for i, j := 0, 1; j < len(t.Tour); i, j = i+1, j+1 {
		node1 = t.Tour[i].NodeNumber
		node2 = t.Tour[j].NodeNumber
		cost += t.DistanceMatrix[node1][node2]
	}
	return cost
}

func generate_matrix(size int, maxDistance int, symmetric bool) [][]int {
	rand.Seed(time.Now().UnixNano())
	// Populate matrix with random values
	grid := make([][]int, size)
	for i := 0; i < size; i++ {
		grid[i] = make([]int, size)
		for j := 0; j < size; j++ {
			grid[i][j] = rand.Intn(maxDistance)
		}
		grid[i][i] = 0
	}

	if symmetric == true {
		// Make matrix symmetrick
		for i := 0; i < size; i++ {
			for j := 0; j < size; j++ {
				grid[i][j] = grid[j][i]
			}
		}
	}

	return grid
}

func (t *TSP) swaps() {
	bestObjective := t.calculate_objective()
	tourLength := t.NumberOfNodes

	for i := 0; i < tourLength; i++ {
		for j := 0; j < tourLength; j++ {
			node1 := t.Tour[i]
			node2 := t.Tour[j]

			t.Tour[i] = node2
			t.Tour[j] = node1

			sampleObjective := t.calculate_objective()

			if sampleObjective < bestObjective {
				bestObjective = sampleObjective
			} else {
				t.Tour[i] = node1
				t.Tour[j] = node2
			}

		}
	}
}

func (t *TSP) performSwaps() {
	lastObjective := t.calculate_objective()
	currentObjective := t.calculate_objective()
	for {
		t.swaps()
		currentObjective = t.calculate_objective()
		if currentObjective == lastObjective {
			break
		} else {
			lastObjective = currentObjective
		}
	}
}

func (t *TSP) pprint() {
	nodeNumbers := make([]int, t.NumberOfNodes)
	for i := 0; i < t.NumberOfNodes; i++ {
		nodeNumbers[i] = t.Tour[i].NodeNumber
	}
	fmt.Println(nodeNumbers)
}

func (t *TSP) shuffle() {
	rand.Seed(time.Now().UnixNano())
	newIndices := rand.Perm(t.NumberOfNodes)
	newTour := make([]Node, t.NumberOfNodes)
	for i := 0; i < t.NumberOfNodes; i++ {
		n := Node{NodeNumber: i}
		newTour[newIndices[i]] = n
	}
	t.Tour = newTour
}

func (t *TSP) solve() {
	t.shuffle()
	t.performSwaps()
}

func buildTSPInstance(num int, distMat [][]int) TSP {
	nodeCount := num
	newTour := make([]Node, nodeCount)

	for i := 0; i < nodeCount; i++ {
		n := Node{NodeNumber: i}
		newTour[i] = n
	}

	// distMat := generate_matrix(nodeCount, 1000, true)
	tspInstance := TSP{Tour: newTour,
		DistanceMatrix: distMat,
		NumberOfNodes:  nodeCount}

	tspInstance.solve()

	return tspInstance
}

func main() {
	cityCount := 500
	distMat := generate_matrix(cityCount, 1000, true)
	start := time.Now()
	randomRestartsRun := 10
	solutions := make([]TSP, randomRestartsRun)

	var wg sync.WaitGroup
	for i := 0; i < randomRestartsRun; i++ {
		wg.Add(1)
		go func(j int) {
			defer wg.Done()
			output_tsp := buildTSPInstance(cityCount, distMat)
			solutions[j] = output_tsp
		}(i)
	}
	wg.Wait()
	// Find minimum

	bestSolution := solutions[0]
	bestObjective := bestSolution.calculate_objective()
	for i := 0; i < randomRestartsRun; i++ {
		if solutions[i].calculate_objective() < bestObjective {
			bestSolution = solutions[i]
			bestObjective = bestSolution.calculate_objective()
		}
	}

	fmt.Printf("Best Solution After %v Iterations: %v\n",
		randomRestartsRun, bestObjective)
	elapsed := time.Since(start)
	fmt.Println("Two-Opt took: ", elapsed)

}
