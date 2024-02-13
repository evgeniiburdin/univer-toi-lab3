package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"strconv"
)

// A constant abstract value to transform given rib weight(distance between two nodes) into
// two nodes' "closeness" coefficient. Can be absolutely anything
const CloseCoef float64 = 200

// Used in pheromone calculation formula: pheromone to add to rib = Q / total distance passed by an ant
// Also can be pretty anything
const Q float64 = 4

// Coefficient that we decrease pheromone on all the ribs by with each global iteration.
// Recommended not to change
const PheroCoff float64 = 0.67

// Amount of nodes you will use if your graph
const nodeAmount int = 7

// Amount of global iterations
const iterationAmount int = 10

// Custom type for storing graph nodes with their coordinates.
// "Visited" field provided for preventing an ant to visit nodes which
// it has already visited through the single ant iteration
type Node struct {
	X       float64
	Y       float64
	Visited bool
}

// Custom type used as *value* value of [][]Graph.
// The two fields represent closeness and pheromone amount between two nodes
type Ridge struct {
	Closeness float64
	Phero     float64
}

// Custom type used to convert "desires to go to each node" into intervals
// that fit into [0;1] global interval
type Interval struct {
	LeftBound  float64
	RightBound float64
}

// Main Function
func main() {

	// First we create a map for our graph nodes
	// You can use floats instead of ints, but you'll need to make some little changes in the whole code
	// by yourself
	nodeMap := make(map[int]*Node)

	// VIP: Creating nodes by hardcoding(dont forget to change "nodeAmount" constant)

	nodeMap[0] = &Node{2, 2, false}
	nodeMap[1] = &Node{9, 13, false}
	nodeMap[2] = &Node{2, 9, false}
	nodeMap[3] = &Node{9, 5, false}
	nodeMap[4] = &Node{6, 4, false}
	nodeMap[5] = &Node{4, 14, false}
	nodeMap[6] = &Node{11, 10, false}

	// VIP: Creating nodes by user input
	/*
		for lk := range nodeAmount {
			fmt.Printf("Node %v:\n", lk)
			var tempX, tempY float64
			fmt.Print("    X: ")
			fmt.Scanln(&tempX)
			fmt.Print("    Y: ")
			fmt.Scanln(&tempY)
			tempNode := Node{tempX, tempY, false}
			nodeMap[lk] = &tempNode
		}
	*/

	// Creating a graph to store distances and pheromones between each graph node
	var Graph [nodeAmount][nodeAmount]Ridge

	// Filling our graph with information based on nodes we've taken on the previous step
	for i := range nodeMap {
		for j := 0; j < len(nodeMap); j++ {
			tempRidge := Ridge{}
			if i == j {
				tempRidge.Closeness = 1
			} else {
				tempRidge.Closeness = CloseCoef / (math.Sqrt((math.Abs(nodeMap[i].X-nodeMap[j].X))*(math.Abs(nodeMap[i].X-nodeMap[j].X)) + (math.Abs(nodeMap[i].Y-nodeMap[j].Y))*(math.Abs(nodeMap[i].Y-nodeMap[j].Y))))
			}
			tempRidge.Phero = 0.2
			Graph[i][j] = tempRidge

		}
	}

	// A map that by the end of the global iteration loop will contain info about all the routes that all the ants have passed
	routeMap := make(map[int]string)

	// Global iteration loop
	for kl := range iterationAmount {
		fmt.Println("iteration: ", kl)

		// With each iteration Every ant makes a loop from the starting node
		for node := range nodeMap {

			currentNode := node
			startingNode := currentNode

			for s := range nodeMap {
				nodeMap[s].Visited = false
			}

			var SingleAntRoute string = ""

			// With each iteration A single ant makes a loop
			for {

				// A map which will contain probabilities of choosing each node as the next, basing on closeness and pheromones
				routeProbMap := make(map[int]float64)

				// Filling our map
				for k := range nodeMap {
					if !nodeMap[k].Visited {
						var selectedNodeProb float64 = (Graph[currentNode][k].Phero * Graph[currentNode][k].Closeness)
						var leastNodesProb float64 = 0.0
						for b := 0; b < len(nodeMap); b++ {
							if b != currentNode && !nodeMap[b].Visited {
								leastNodesProb += Graph[currentNode][b].Phero * Graph[currentNode][b].Closeness

							}
						}
						var probability float64
						if currentNode != k {
							probability = ((selectedNodeProb) / (leastNodesProb))
						} else {
							probability = 0
						}

						routeProbMap[k] = probability
					}
				}

				// Converting our map to intervals on the [0;1] line segment

				intervalMap := make(map[int]Interval)
				var LeftBound float64 = 0
				for l := range routeProbMap {

					if routeProbMap[l] != 0 {
						tempInterval := Interval{LeftBound, LeftBound + routeProbMap[l]}
						intervalMap[l] = tempInterval
						LeftBound = tempInterval.RightBound
					}
				}

				// *Randomly* choosing the way the ant will go basing on each probability interval scale

				var randomVar float64 = rand.Float64()

				var selectedNode int
				for h := range intervalMap {
					if randomVar >= intervalMap[h].LeftBound && randomVar < intervalMap[h].RightBound {
						selectedNode = h
						break
					}
				}

				nodeMap[currentNode].Visited = true
				SingleAntRoute += strconv.Itoa(currentNode)
				currentNode = selectedNode

				// Checking if all the nodes have already been visited.
				var visitedAll = true
				for l := range nodeMap {
					if !nodeMap[l].Visited {
						visitedAll = false
					}
				}

				// If yes - breaking the infinite loop
				if visitedAll {
					break
				}
			}

			// Ending out single ant route with the starting point, as the ant making a loop across all the nodes
			SingleAntRoute += strconv.Itoa(startingNode)

			// Passing the route that a single ant passed to a map containing all the routes of all the ants
			routeMap[startingNode] = SingleAntRoute

		}

		// Calculating the overall distance the single ant has passed
		var initPoint, nextPoint int
		var err error
		var totalDist float64
		for n := 0; n < len(routeMap); n++ {

			currentRoute := routeMap[n]

			totalDist = 0
			for o := 0; o < len(currentRoute)-1; o++ {

				initPoint, err = strconv.Atoi(string(currentRoute[o]))

				if err != nil {
					log.Fatal(err)
				}
				nextPoint, err = strconv.Atoi(string(currentRoute[o+1]))

				if err != nil {
					log.Fatal(err)
				}

				totalDist += CloseCoef / Graph[initPoint][nextPoint].Closeness

			}

		}

		// Evaporation of pheromones
		for i := 0; i < nodeAmount; i++ {
			for j := 0; j < nodeAmount; j++ {
				Graph[i][j].Phero *= PheroCoff
			}
		}

		// Updating pheromones depending on which ribs the ants have visited most of all
		for n := 0; n < len(routeMap); n++ {
			currentRoute := routeMap[n]
			var pheromone float64 = Q / totalDist
			for o := 0; o < len(currentRoute)-1; o++ {
				initPoint, err = strconv.Atoi(string(currentRoute[o]))
				if err != nil {
					log.Fatal(err)
				}
				nextPoint, err = strconv.Atoi(string(currentRoute[o+1]))
				if err != nil {
					log.Fatal(err)
				}
				Graph[initPoint][nextPoint].Phero += pheromone

			}
		}

		// Printing out out Graph with each iteration

		fmt.Print("\n\n\n")
		for f := 0; f < nodeAmount; f++ {
			//fmt.Print("            ")
			for u := 0; u < nodeAmount; u++ {
				fmt.Print("D: ", fmt.Sprintf("%.2f", CloseCoef/Graph[f][u].Closeness), "    ")
			}
			fmt.Println()
			for u := 0; u < nodeAmount; u++ {
				fmt.Print("Ph:", fmt.Sprintf("%.2f", Graph[f][u].Phero), "    ")
			}

			fmt.Println()
			fmt.Println()
		}
		fmt.Print("\n\n\n")

	}

	// In a couple of iterations from the starting point we might want to know the optimal way across
	// all the nodes in our graph.
	// This part is simulating a single ant choosing the next node each time depending only on pheromones
	currentNode := 0
	for node := range nodeMap {
		nodeMap[node].Visited = false
	}

	var route string = ""

	for {

		var bestWay int
		var CurrentPhero float64 = 0

		for gh := range nodeMap {
			if gh != currentNode && !nodeMap[gh].Visited {
				if Graph[currentNode][gh].Phero > CurrentPhero {
					CurrentPhero = Graph[currentNode][gh].Phero
					bestWay = gh
				}
			}
		}

		nodeMap[currentNode].Visited = true
		route += strconv.Itoa(currentNode + 1)
		currentNode = bestWay

		var visitedAll = true
		for l := range nodeMap {
			if !nodeMap[l].Visited {
				visitedAll = false
			}
		}
		if visitedAll {
			break
		}
	}

	fmt.Println("OPTIMAL WAY: ", route)

}
