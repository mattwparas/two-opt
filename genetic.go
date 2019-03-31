package main

import (
	"fmt"
	"math/rand"
	"time"
	"sort"
)

type Node struct {
	NodeNumber int
}

type TSP struct {
	Tour           []Node
	DistanceMatrix [][]float64
	NumberOfNodes  int
	RouteDistance float64
	Fitness float64
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (t *TSP) calculate_fitness() float64 {
	if t.Fitness == 0 {
		t.Fitness = 1 / t.RouteDistance
	} else {
		t.Fitness = 1 / t.RouteDistance
	}
	return t.Fitness
}

func (t *TSP) calculate_objective() float64 {
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
	t.RouteDistance = cost
	return cost
}

func generate_matrix(size int, maxDistance float64, symmetric bool) [][]float64 {
	rand.Seed(27)
	// Populate matrix with random values
	grid := make([][]float64, size)
	for i := 0; i < size; i++ {
		grid[i] = make([]float64, size)
		for j := 0; j < size; j++ {
			grid[i][j] = maxDistance * rand.Float64()
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
	// rand.Seed(time.Now().UnixNano())
	newIndices := rand.Perm(t.NumberOfNodes)
	newTour := make([]Node, t.NumberOfNodes)
	for i := 0; i < t.NumberOfNodes; i++ {
		n := Node{NodeNumber: i}
		newTour[newIndices[i]] = n
	}
	t.Tour = newTour
}

func (t *TSP) solve() {
	// t.shuffle()
	t.performSwaps()
	t.RouteDistance = t.calculate_objective()
}

func buildTSPInstance(num int, distMat [][]float64) TSP {
	nodeCount := num
	newTour := make([]Node, nodeCount)

	for i := 0; i < nodeCount; i++ {
		n := Node{NodeNumber: i}
		newTour[i] = n
	}

	// distMat := generate_matrix(nodeCount, 1000, true)
	tspInstance := TSP{Tour: newTour,
		DistanceMatrix: distMat,
		NumberOfNodes:  nodeCount,
		RouteDistance: 0,
		Fitness: 0}

	// Shuffle for initial fitness
	tspInstance.shuffle()

	tspInstance.calculate_objective()
	tspInstance.calculate_fitness()

	// tspInstance.solve()

	return tspInstance
}

func initialPopulation(cityCount int, popSize int) []TSP {
	distMat := generate_matrix(cityCount, 300, true)
	initialPop := make([]TSP, popSize)
	for i := 0; i < popSize; i++ {
		rand.Seed(int64(i))
		initialPop[i] = buildTSPInstance(cityCount, distMat)
	}
	return initialPop
}


func rankRoutes(population []TSP) []TSP {
	sort.Slice(population, func(i, j int) bool {
		return population[i].Fitness > population[j].Fitness
	})
	return population
}

// TODO
func matingPool(population []TSP, eliteSize int) []TSP {

	var matingPool []TSP

	cumSum := make([]float64, len(population))
	cumPerc := make([]float64, len(population))
	cumSum[0] = 1 / float64(population[0].RouteDistance)
	totalSum := 1 / float64(population[0].RouteDistance)

	// calculate the cumulative sum of fitness
	for i:=1; i < len(population); i++ {
		totalSum += 1.0 / float64(population[0].RouteDistance)
		cumSum[i] = totalSum
	}

	for i:=0; i < len(population); i++ {
		cumPerc[i] = float64(cumSum[i]) / float64(totalSum)
		// fmt.Println(cumSum[i])
		// fmt.Println(cumPerc[i])
	}

	matingPool = append(matingPool, population[0:eliteSize]...)

	// Double Picking here is bad (maybe?)
	for i:=0; i < len(population) - eliteSize; i++ {
		pick := rand.Float64()
		for j:=0; j < len(population); j++ {
			if pick <= cumPerc[j] {
				matingPool = append(matingPool, population[j])
				break
			}
		}
	}

	// for i:=0; i < len(population) - eliteSize; i++ {
	// 	matingPool = append(matingPool, population[rand.Intn(len(population) - eliteSize)])
	// }


	return matingPool
}

func breed(parent1 TSP, parent2 TSP) TSP {
	tourLength := parent1.NumberOfNodes
	geneA := rand.Intn(tourLength)
	geneB := rand.Intn(tourLength)

	startGene := min(geneA, geneB)
	endGene := max(geneA, geneB)

	var child1Tour []Node
	var child2Tour []Node
	//remainTour := make([]Node, tourLength)
	nodeMap := make(map[int]bool)

	for i := startGene; i < endGene; i++ {
		child1Tour = append(child1Tour, parent1.Tour[i])
		nodeMap[parent1.Tour[i].NodeNumber] = true
	}

	// fill in the gaps with everything else
	for i := 0; i < tourLength; i++ {
		if nodeMap[parent2.Tour[i].NodeNumber] == false {
			// child.Tour[i] = parent2.Tour[i]
			child2Tour = append(child2Tour, parent2.Tour[i])
		}
	}

	child := TSP{Tour: append(child1Tour, child2Tour...),
		DistanceMatrix: parent1.DistanceMatrix,
		NumberOfNodes:  tourLength,
		RouteDistance: 0,
		Fitness: 0}

	child.calculate_objective()
	child.calculate_fitness()

	return child
}

func breedPopulation(matingPool []TSP, eliteSize int) []TSP {
	var children []TSP
	length := len(matingPool) - eliteSize
	//rand.Seed(time.Now().UnixNano())
	newIndices := rand.Perm(len(matingPool))
	newPool := make([]TSP, len(matingPool))
	for i := 0; i < len(newPool); i++ {
		n := matingPool[i]
		newPool[newIndices[i]] = n
	}
	children = append(children, matingPool[0:eliteSize]...)
	for i := 0; i < length; i++ {
		child := breed(newPool[i], newPool[len(matingPool) - i - 1])
		children = append(children, child)
	}
	return children
}

func (individual *TSP) mutate(mutationRate float64) {
	for i := 0; i < individual.NumberOfNodes; i++ {
		if rand.Float64() < mutationRate{
			j := rand.Intn(individual.NumberOfNodes)
			
			node1 := individual.Tour[i]
			node2 := individual.Tour[j]

			individual.Tour[i] = node2
			individual.Tour[j] = node1
		}
	}
}

func mutatePopulation(population []TSP, mutationRate float64) []TSP {
	for i := 0; i < len(population); i++ {
		// population[i].mutate(mutationRate)
		if rand.Float64() < mutationRate {
			population[i].solve()
		}
		// population[i].solve()
	}
	return population
}

func nextGeneration(currentGen []TSP, eliteSize int, mutationRate float64) []TSP {
	popRanked := rankRoutes(currentGen)
	matingPool := matingPool(popRanked, eliteSize)
	children := breedPopulation(matingPool, eliteSize)
	nextGeneration := mutatePopulation(children, mutationRate)
	return nextGeneration
}

func geneticAlgorithm(popSize int, eliteSize int, mutationRate float64, generations int) TSP {
	rand.Seed(27)
	pop := initialPopulation(100, popSize)
	bestInitialDistance := rankRoutes(pop)[0].RouteDistance
	// initialBest := rankRoutes(pop)[0]
	// initialBest.solve()
	fmt.Println("Initial Distance: ", bestInitialDistance)
	// fmt.Println("Initial best 2 opt: ", initialBest.RouteDistance)

	for i := 0; i < generations; i++ {
		// debugGeneration(pop)
		pop = nextGeneration(pop, eliteSize, mutationRate)
		// fmt.Println(pop[0].RouteDistance)
		// fmt.Println(i)
	}

	bestTSP := rankRoutes(pop)[0]
	// fmt.Println("Best Distance before 2 opt: ", bestTSP.RouteDistance)
	// bestTSP.solve()
	fmt.Println("Final Distance: ", bestTSP.RouteDistance)
	// bestTSP.pprint()
	return bestTSP
}

func main() {
	start := time.Now()
	geneticAlgorithm(1000, 200, 0.01, 150)
	elapsed := time.Since(start)
	fmt.Println("Genetic Algorithm took: ", elapsed)
}














