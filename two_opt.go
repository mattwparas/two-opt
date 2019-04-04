package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// represents a city
type Node struct {
	NodeNumber int
}

// represents a wholistic solution - a "route"
type TSP struct {
	Tour           []Node
	DistanceMatrix [][]int
	NumberOfNodes  int
}

func (t *TSP) calculateObjective() int {
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

func generateMatrix(size int, maxDistance int, symmetric bool) [][]int {
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

func (t *TSP) marginalCost(i int) int {
	nodeBefore := i - 1
	node := i
	nodeAfter := (i + 1) % t.NumberOfNodes
	if i == 0 {
		nodeBefore = t.NumberOfNodes - 1
	}
	marginalCost := t.DistanceMatrix[t.Tour[nodeBefore].NodeNumber][t.Tour[node].NodeNumber]
	marginalCost += t.DistanceMatrix[t.Tour[node].NodeNumber][t.Tour[nodeAfter].NodeNumber]
	return marginalCost
}

func (t *TSP) swaps() {
	// bestObjective := t.calculateObjective()
	tourLength := t.NumberOfNodes

	for i := 0; i < tourLength; i++ {
		for j := 0; j < tourLength; j++ {
			node1 := t.Tour[i]
			node2 := t.Tour[j]

			// cost from i-1 to i and i to i + 1
			// cost from j-1 to j and j to j + 1
			currentCost := t.marginalCost(i) + t.marginalCost(j)

			t.Tour[i] = node2
			t.Tour[j] = node1

			newCost := t.marginalCost(i) + t.marginalCost(j)
			if newCost > currentCost {
				t.Tour[i] = node1
				t.Tour[j] = node2
			}
		}
	}
}

func (t *TSP) performSwaps() {
	lastObjective := t.calculateObjective()
	currentObjective := t.calculateObjective()
	for {
		t.swaps()
		currentObjective = t.calculateObjective()
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

	// distMat := generateMatrix(nodeCount, 1000, true)
	tspInstance := TSP{Tour: newTour,
		DistanceMatrix: distMat,
		NumberOfNodes:  nodeCount}

	tspInstance.solve()

	return tspInstance
}

func main() {
	cityCount := 1500
	distMat := generateMatrix(cityCount, 1000, true)
	start := time.Now()
	randomRestartsRun := 10
	solutions := make([]TSP, randomRestartsRun)

	var wg sync.WaitGroup
	for i := 0; i < randomRestartsRun; i++ {
		wg.Add(1)
		go func(j int) {
			defer wg.Done()
			outputTsp := buildTSPInstance(cityCount, distMat)
			solutions[j] = outputTsp
		}(i)
	}
	wg.Wait()
	// Find minimum

	bestSolution := solutions[0]
	bestObjective := bestSolution.calculateObjective()
	for i := 0; i < randomRestartsRun; i++ {
		if solutions[i].calculateObjective() < bestObjective {
			bestSolution = solutions[i]
			bestObjective = bestSolution.calculateObjective()
		}
	}

	fmt.Printf("Best Solution After %v Iterations: %v\n",
		randomRestartsRun, bestObjective)
	elapsed := time.Since(start)
	fmt.Println("Two-Opt took: ", elapsed)

}
