package main

import (
	"fmt"
	"math/rand"
	"strings"
	"os"
	"os/exec"
	"bufio"
	"strconv"
	"time"
)

// TODO: Check if puzzle is solveable, http://www.cs.bham.ac.uk/~mdr/teaching/modules04/java2/TilesSolvability.html
// TODO: Check os before clearing screen

var scanner *bufio.Scanner

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
	scanner = bufio.NewScanner(os.Stdin)
}

type PuzzlePiece struct {
	Value string
}

type Puzzle struct {
	Width int
	Height int
	Pieces []PuzzlePiece
	EmptyIndex int
}

func (puzzle *Puzzle) String() string {
	var b strings.Builder
	
	for i, v := range puzzle.Pieces {
		if i > 0 && i % puzzle.Width == 0 {
			b.WriteString("\n")
		}
		b.WriteString(v.Value)
		b.WriteString(" ")
	}
	
	return b.String()
}

func (puzzle *Puzzle) Generate(width, height int) {
	puzzle.Width = width
	puzzle.Height = height
	puzzle.Pieces = make([]PuzzlePiece, width * height)
	
	randomOrder := rand.Perm(len(puzzle.Pieces))
	for i, v := range randomOrder {
		if v == len(puzzle.Pieces) - 1 {
			puzzle.Pieces[i] = PuzzlePiece{ "__" }
			puzzle.EmptyIndex = i
		} else {
			puzzle.Pieces[i] = PuzzlePiece{ fmt.Sprintf("%02d", v + 1) }
		}
	}
}

func (puzzle *Puzzle) TranslateInput(input string) (relativeIndex int) {
	switch input {
		case "a":
			relativeIndex = -1
		case "d":
			relativeIndex = 1
		case "w":
			relativeIndex = -puzzle.Width
		case "s":
			relativeIndex = puzzle.Width
		default:
			relativeIndex = 0
	}
	
	return
}

func (puzzle *Puzzle) HandleInput(input string) {
	relativeIndex := puzzle.TranslateInput(input)
	newIndex := puzzle.EmptyIndex + relativeIndex
	
	// Out of bounds
	if newIndex < 0 || newIndex >= len(puzzle.Pieces) {
		return
	}
	
	// Check if newIndex is on the same column or row as current index
	cX, cY := puzzle.EmptyIndex % puzzle.Width, puzzle.EmptyIndex / puzzle.Height
	rX, rY := newIndex % puzzle.Width, newIndex / puzzle.Height
	
	if cX == rX || cY == rY {
		puzzle.Pieces[puzzle.EmptyIndex], puzzle.Pieces[newIndex] = puzzle.Pieces[newIndex], puzzle.Pieces[puzzle.EmptyIndex]
		puzzle.EmptyIndex = newIndex
	}
}

func GetUserInput() (input string) {
	fmt.Println("")
	fmt.Println("Select direction using w, a, s, d. (type exit to end)")
	fmt.Print(">")
	
	scanner.Scan()
	input = scanner.Text()
	return
}

func (puzzle *Puzzle) Won() bool {
	// Puzzle can't be won if the last piece isn't empty
	if puzzle.EmptyIndex != len(puzzle.Pieces) - 1 {
		return false
	}

	for i, v := range puzzle.Pieces {
		if v.Value == "__" {
			continue
		}
		
		num, _ := strconv.Atoi(v.Value)
		if num - i != 1 {
			return false
		}
	}
	
	return true
}

func (puzzle *Puzzle) Play() {
	for {
		// Show state
		fmt.Println(puzzle.String())
	
		// Get user input
		input := GetUserInput()
		
		// Act on input
		if input == "exit" {
			fmt.Println("Game exited by player")
			break
		} else {
			puzzle.HandleInput(input)
		}
		
		ClearWindows()
		
		// Check for win
		if puzzle.Won() {
			fmt.Println(puzzle.String())
			fmt.Println("Puzzle complete!")
			break
		} 
	}
}

func GetPuzzleFromUser() *Puzzle {
	fmt.Println("")
	fmt.Println("Sliding puzzle game!")
	fmt.Print("Size: ")
	scanner.Scan()
	fmt.Println("")
	input := scanner.Text()
	
	size, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("Please enter a number, error:")
		fmt.Println(err)
		return nil
	} else {
		defer ClearWindows()
	}
	
	puzzle := Puzzle{}
	puzzle.Generate(size, size)
	
	return &puzzle
}

func main() {
	if puzzle := GetPuzzleFromUser(); puzzle != nil {
		puzzle.Play()
	}
}

func ClearWindows() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}