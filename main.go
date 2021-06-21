package main

import (
	"fmt"
	"github.com/draffensperger/golp"
	"log"
)

const (
	numDays = 30
)

type Constraints struct {
	StretchLimit int
	NumShifts    float64
	IsAvailable  func(day int) bool
}

var (
	// (person i) is available on these given shifts
	constraintsByPerson = map[int]Constraints{
		0: {
			StretchLimit: 3,
			NumShifts:    7,
			IsAvailable: func(day int) bool {
				return day < 14
			},
		},
		1: {
			StretchLimit: 4,
			NumShifts:    8,
			IsAvailable: func(day int) bool {
				return day > 7 && day < 21
			},
		},
		2: {
			StretchLimit: 5,
			NumShifts:    7,
			IsAvailable: func(day int) bool {
				return day > 14
			},
		},
		3: {
			StretchLimit: 2,
			NumShifts:    8,
			IsAvailable: func(day int) bool {
				return true
			},
		},
	}
	numPeople = len(constraintsByPerson)
	numVars   = numPeople * numDays
)

func main() {
	lp := golp.NewLP(0, numVars)

	// ensure no one is scheduled when they are not available
	for person, constraints := range constraintsByPerson {
		// ensure they don't have stretches that are too long
		{
			for start := 0; start+constraints.StretchLimit < numDays; start++ {
				row := make([]float64, numVars)
				for x := 0; x < constraints.StretchLimit; x++ {
					row[getIndex(person, start+x)] = 1
				}
				if err := lp.AddConstraint(row, golp.LE, float64(constraints.StretchLimit)-1); err != nil {
					log.Fatal(err)
				}
			}
		}

		// ensure they have the right number of shifts
		{
			row := make([]float64, numVars)
			for day := 0; day < numDays; day++ {
				row[getIndex(person, day)] = 1
			}
			if err := lp.AddConstraint(row, golp.EQ, constraints.NumShifts); err != nil {
				log.Fatal(err)
			}
		}

		// ensure they are only scheduled when available
		{
			row := make([]float64, numVars)
			for day := 0; day < numDays; day++ {
				val := float64(1)
				if constraints.IsAvailable(day) {
					val = 0
				}
				row[getIndex(person, day)] = val
			}
			if err := lp.AddConstraint(row, golp.EQ, 0); err != nil {
				log.Fatal(err)
			}
		}
	}

	// ensure one person is scheduled at all times
	for day := 0; day < numDays; day++ {
		row := make([]float64, numVars)
		for person := 0; person < numPeople; person++ {
			row[getIndex(person, day)] = 1
		}
		if err := lp.AddConstraint(row, golp.EQ, 1); err != nil {
			log.Fatal(err)
		}
	}

	// force integers
	for person := 0; person < numPeople; person++ {
		for day := 0; day < numDays; day++ {
			lp.SetInt(getIndex(person, day), true)
		}
	}

	lp.SetObjFn(make([]float64, numVars))
	fmt.Println(lp.Solve())
	for day := 0; day < numDays; day++ {
		for person := 0; person < numPeople; person++ {
			if lp.Variables()[getIndex(person, day)] > 0 {
				fmt.Printf("Day %d --> Person %d\n", day, person)
			}
		}
	}
}

func getIndex(person int, day int) int {
	return person*numDays + day
}
