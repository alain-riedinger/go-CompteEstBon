package main

// Result contains the information needed for the recursive solving
type Result struct {
	Value int
	Text  string
	Steps []int
}

// NewResultPtr initializes a Result structure, with pointer
func NewResultPtr() *Result {
	r := new(Result)
	return r
}

// NewResult initializes a Result structure
func NewResult() Result {
	r := new(Result)
	return *r
}
