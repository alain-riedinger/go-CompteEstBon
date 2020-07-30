package main

// Solution contains the information needed for the recursive solving
type Solution struct {
	Tirage  int
	Best    Result
	Depth   int
	Current []Result
}

// NewSolution initializes a Solution structure
func NewSolution() *Solution {
	s := new(Solution)
	return s
}
