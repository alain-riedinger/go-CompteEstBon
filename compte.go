package main

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"time"
)

// NbTirage is the number of plaques in the tirage
const nbTirage = 6

// Games is composed on 24 tiles, with following occurences:
//   1 to 10 = twice
//   25, 50, 75, 100 = once
const nbPlaques = 24
const minTirage = 101
const maxTirage = 999

// Compte is the class that holds the game
type Compte struct {
	plaques []int
}

// NewCompte initializes the Compte structure
func NewCompte() *Compte {
	c := new(Compte)
	c.plaques = make([]int, nbPlaques)

	idx := 0
	// 1 to 10 is present twice
	for n := 1; n <= 10; n++ {
		c.plaques[idx] = n
		c.plaques[idx+1] = n
		idx += 2
	}
	// Other numbers are present once
	c.plaques[idx] = 25
	idx++
	c.plaques[idx] = 50
	idx++
	c.plaques[idx] = 75
	idx++
	c.plaques[idx] = 100
	idx++

	// Initializes uniquely the random
	rand.Seed(time.Now().UnixNano())

	// Lets' shuffle a bit the seeds to avoid artefacts due to random
	c.plaques = shuffle(c.plaques)

	return c
}

// GetPlaques returns a random tirage
// with nbPlaques numbers
func (c *Compte) GetPlaques() []int {
	var chosenIdx []int

	// Randomly select a list of indexes (can only be chosen once)
	for nbChosen := 0; nbChosen < nbTirage; {
		p := rand.Intn(nbPlaques - 1)
		if !contains(chosenIdx, p) {
			chosenIdx = append(chosenIdx, p)
			nbChosen++
		}
	}

	// Compute the returned list with the values
	tirage := make([]int, nbTirage)
	for i := 0; i < nbTirage; i++ {
		tirage[i] = c.plaques[chosenIdx[i]]
	}
	return tirage
}

// GetTirage returns a random number to calculate
func (c *Compte) GetTirage() int {
	return minTirage + rand.Intn(maxTirage-minTirage)
}

// SolveTirage recursively solves a computation of a tirage
// - recursive
// - apply the 4 operations
// - stop if found the result
// - if not found, keep the closest result
// - results must be chained: it can approach, then go farther, and closer after
// - first keep the first correct result
// - in a more sophisticated version, store all good results and propose shorest
//
// Principles of operations:
// - "+" and "x" can always be applied
// - "-" can only produce strictly positive values
// - "/" can only produce entire values and divide by 1 is useless
//
// Structural principles:
// - impossible to compute a step with itself
//
// Returns the solution
func (c *Compte) SolveTirage(solution Solution) *Solution {
	// Terminal case: nothing left to calculate
	if solution.Depth == 1 {
		return &solution
	}

	// Loop through all the possible operations for this set of steps
	// Constructs a grid, with only the sub left half to be processed
	// - numbers are sorted ascending
	// - useless to do adition twice: l + r, and then r + l
	// - same with multiply
	// - substraction and division are only one way and left is always bigger
	// Example for a steps set with 6 items:
	//   | 1 | 2 | 3 | 4 | 5 | 6 |
	// 1 | W | U | U | U | U | U |
	// 2 | + | W | U | U | U | U |
	// 3 | + | + | W | U | U | U |
	// 4 | + | + | + | W | U | U |
	// 5 | + | + | + | + | W | U |
	// 6 | + | + | + | + | + | W |
	var newCur []Result
	for _, cur := range solution.Current {
		for l := 1; l < len(cur.Steps); l++ {
			for r := 0; r < l; r++ {
				// Loop through the 4 possible operations
				for o := 0; o < 4; o++ {
					res := calculate(cur.Steps, l, r, o)
					if res != nil {
						res.Text = cur.Text + res.Text + "\n"

						if math.Abs(float64(res.Value-solution.Tirage)) < math.Abs(float64(solution.Best.Value-solution.Tirage)) {
							solution.Best = *res
						}

						// Found the exact solution: terminal return
						if solution.Best.Value == solution.Tirage {
							return &solution
						}
						newCur = append(newCur, *res)
					}
				}
			}
		}
	}

	// Update the list of possibilities with the new calculated ones
	solution.Current = newCur
	// Decrease depth of next level of steps
	solution.Depth--

	// Recursive call for next level
	return c.SolveTirage(solution)
}

func calculate(steps []int, l int, r int, op int) *Result {
	result := NewResult()
	// Allocate a new set of steps that is 1 step smaller
	result.Steps = make([]int, len(steps)-1)
	// Copy the non modified steps (not in calculation)
	n := 0
	for i := 0; i < len(steps); i++ {
		if i != l && i != r {
			result.Steps[n] = steps[i]
			n++
		}
	}
	// Perform the calculation, if it has a sense
	// Senseful returns are done inside of the switch, useless null return done at the end of the method
	switch op {
	case 0: // Add
		result.Steps[n] = steps[l] + steps[r]
		result.Value = result.Steps[n]
		result.Text = fmt.Sprintf("%d + %d = %d\n", steps[l], steps[r], result.Value)
		return &result
	case 1: // Multiply
		if steps[r] != 1 {
			result.Steps[n] = steps[l] * steps[r]
			result.Value = result.Steps[n]
			result.Text = fmt.Sprintf("%d x %d = %d\n", steps[l], steps[r], result.Value)
			return &result
		}
		// A multiply by 1 is no help at all: skip this possibility
	case 2: // Substract
		if steps[l] > steps[r] {
			result.Steps[n] = steps[l] - steps[r]
			result.Value = result.Steps[n]
			result.Text = fmt.Sprintf("%d - %d = %d\n", steps[l], steps[r], result.Value)
			return &result
		}
		// A substract that leads to a value of 0 is no help at all: skip this possibility
	case 3: // Divide
		if steps[r] != 1 && steps[l]%steps[r] == 0 {
			result.Steps[n] = steps[l] / steps[r]
			result.Value = result.Steps[n]
			result.Text = fmt.Sprintf("%d / %d = %d\n", steps[l], steps[r], result.Value)
			return &result
		}
		// A divide that is not entire is not possible
		// A divide by 1 is no help at all: skip these possibilities
	}
	// When here, means that calculation is not possible or leads to nothing helpful
	return nil
}

func contains(s interface{}, elem interface{}) bool {
	arrV := reflect.ValueOf(s)
	if arrV.Kind() == reflect.Slice {
		for i := 0; i < arrV.Len(); i++ {
			// XXX - panics if slice element points to an unexported struct field
			// see https://golang.org/pkg/reflect/#Value.Interface
			if arrV.Index(i).Interface() == elem {
				return true
			}
		}
	}
	return false
}

func shuffle(numbers []int) []int {
	rand.Shuffle(len(numbers), func(i, j int) {
		numbers[i], numbers[j] = numbers[j], numbers[i]
	})
	return numbers
}
