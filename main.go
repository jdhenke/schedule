package main

import (
	"fmt"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/optimize/convex/lp"
	"log"
)

const (
	numDays = 7
)

var (
	// (person i) is available on these given shifts
	availabilityByPerson = map[int][]int{
		0: {0, 1, 2, 3},
		1: {4, 5, 6},
		2: {0, 2, 4, 6},
	}
	numPeople = len(availabilityByPerson)
	numVars   = numPeople * numDays
)

func main() {
	// just want *A* solution for now
	c := make([]float64, numVars)

	// will assemble constraints here
	var A []float64
	var b []float64

	// ensure no one is scheduled when they are not available
	for person, availability := range availabilityByPerson {
		// row of A ensuring they are not working when they are NOT available
		row := make([]float64, numVars)
		for day := 0; day < numDays; day++ {
			row[getIndex(person, day)] = 1
		}
		for _, day := range availability {
			row[getIndex(person, day)] = 0
		}
		A = append(A, row...)
		b = append(b, 0)
	}

	// ensure someone is scheduled at all times
	for day := 0; day < numDays; day++ {
		row := make([]float64, numVars)
		for person := 0; person < numPeople; person++ {
			row[getIndex(person, day)] = 1
		}
		A = append(A, row...)
		b = append(b, 1)
	}

	// solve and print solution
	matA := mat.NewDense(len(A)/numVars, numVars, A)
	_, x, err := lp.Simplex(c, matA, b, 0, nil)
	if err != nil {
		log.Fatal(err)
	}
	for day := 0; day < numDays; day++ {
		for person := 0; person < numPeople; person++ {
			if x[getIndex(person, day)] == 1 {
				fmt.Printf("Day %d --> Person %d\n", day, person)
			}
		}
	}
}

func getIndex(person int, day int) int {
	return person*numDays + day
}
