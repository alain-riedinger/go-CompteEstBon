package main

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
)

// ArgError signals an argument error
type ArgError struct {
	arg string
	err string
}

func (e *ArgError) Error() string {
	return fmt.Sprintf("Argument: \"%s\" %s", e.arg, e.err)
}

func help() {
	fmt.Println()
	fmt.Println("+----------------+")
	fmt.Println("| Compte est Bon |")
	fmt.Println("+----------------+")
	fmt.Println("A game that consists in reaching a number with a given set of random numbers")
	fmt.Println("- 6 random numbers among: 1 to 10, 25, 50 75, 100")
	fmt.Println("- the 4 arithmetical operations: +, -, x, /")
	fmt.Println("- calculating results must be entire and positive")
	fmt.Println()
	fmt.Println("Syntax:")
	fmt.Println("-------")
	fmt.Println("CompteEstBon")
	fmt.Println("- to play a game")
	fmt.Println("  no options needed, the application will do")
	fmt.Println()
}

func parseArgs(args []string) error {
	i := 1
	for i < len(args) {
		switch strings.ToLower(args[i]) {
		case "--help", "/help", "-h", "/h":
			help()
			os.Exit(1)
		default:
			return &ArgError{args[i], "unknown"}
		}
	}
	return nil
}

func main() {
	// Parse arguments
	err := parseArgs(os.Args)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		help()
		return
	}

	for {
		clearScreen()
		displayTitle()
		// Launch a game with 40 seconds of think time
		newGame(30, 40)
		prompt("Want another game ?")
	}
}

func findSolution(cpt *Compte, tirage int, plaques []int, sol chan *Solution) {
	// Initialize the recursive calculation root structure
	solution := NewSolution()
	solution.Tirage = tirage
	solution.Depth = len(plaques)
	res := NewResult()
	res.Steps = plaques
	res.Text = ""
	res.Value = 0
	sort.Ints(res.Steps)
	solution.Current = append(solution.Current, res)

	// Initialize the best approaching structure
	solution.Best = NewResult()
	solution.Best.Steps = res.Steps
	solution.Best.Value = solution.Best.Steps[len(solution.Best.Steps)-1]
	solution.Best.Text = fmt.Sprintf("%d", solution.Best.Value)

	// Start the recursive resolution
	solution = cpt.SolveTirage(*solution)

	// Send solution to the channel, while execution in parallel
	sol <- solution
}

func newGame(chars int, seconds int) {
	cpt := NewCompte()
	plaques := cpt.GetPlaques()

	s := " "
	for i := 0; i < len(plaques); i++ {
		if s != "" {
			s += " "
		}
		s += fmt.Sprintf("|%3d|", plaques[i])
	}
	fmt.Println("  +---+ +---+ +---+ +---+ +---+ +---+")
	fmt.Println(s)
	fmt.Println("  +---+ +---+ +---+ +---+ +---+ +---+")

	tirage := cpt.GetTirage()
	t := fmt.Sprintf("  |%3d|", tirage)
	fmt.Println("  +---+")
	fmt.Println(t)
	fmt.Println("  +---+")

	// Wait for some seconds of think time
	countup(chars, seconds)

	// Solution is searched during chrono time (it's shorter so no cheating)
	sol := make(chan *Solution)
	go findSolution(cpt, tirage, plaques, sol)
	solution := <-sol
	close(sol)

	prompt("Want a solution ?")

	// Output final result
	state := "Exact"
	if solution.Best.Value != solution.Tirage {
		state = "Approché"
	}
	fmt.Printf("Solution [%s]\n", state)
	fmt.Println()
	fmt.Println(solution.Best.Text)
	fmt.Println()
}

func displayTitle() {
	fmt.Println("  +-------------------------+")
	fmt.Println("  |  Le Compte est Bon !    |")
	fmt.Println("  +-------------------------+")
}

func prompt(msg string) {
	fmt.Printf("%s (press Enter)\n", msg)
	var s string
	fmt.Scanln(&s)
}

func clearScreen() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
