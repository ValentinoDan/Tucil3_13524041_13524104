package filehandler

import (
	"bufio"
	"fmt"
	"iceSlidingPuzzle/src/puzzle"
	"os"
	"strconv"
	"strings"
)

func parseBoard(filename string) (*puzzle.Board, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("Failed to open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Read first line to get N and M
	if !scanner.Scan() {
		return nil, fmt.Errorf("Empty file / Failed to read first line")
	}
	firstLine := strings.TrimSpace(scanner.Text())

	// read N and M from first line
	var N, M int
	_, err = fmt.Sscanf(firstLine, "%d %d", &N, &M)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse N and M: %w", err)
	}

	// Init board
	board := &puzzle.Board{
		N:        N,
		M:        M,
		Grid:     make([][]rune, N),
		Cost:     make([][]int, N), 
		Obstacle: make(map[int]puzzle.Point),
		Lava:     []puzzle.Point{},
	}

	// read the grid lines
	for row := 0; row < N; row++ {
		if !scanner.Scan() {
			return nil, fmt.Errorf("Expected %d lines but got fewer", N)
		}
		line := strings.TrimSpace(scanner.Text()) // remove space

		// validate line length
		if len(line) != M {
			return nil, fmt.Errorf("Invalid board: row %d has %d columns instead of %d", row+1, len(line), M)
		}

		board.Grid[row] = make([]rune, M)
		for col, char := range line {
			board.Grid[row][col] = char
			switch char {
			case 'Z':
				board.Start = puzzle.Point{Row: row, Col: col}
			case 'O':
				board.Goal = puzzle.Point{Row: row, Col: col}
			case 'L':
				board.Lava = append(board.Lava, puzzle.Point{Row: row, Col: col})
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				board.Obstacle[int(char-'0')] = puzzle.Point{Row: row, Col: col}
			}
		}
	}

	// read the cost lines
	for row := 0; row < N; row++ {
		if !scanner.Scan() {
			return nil, fmt.Errorf("Expected %d lines but got fewer", N)
		}
		line := strings.TrimSpace(scanner.Text()) // remove space

		fields := strings.Fields(line) // parse multipe digit number

		// validate line length
		if len(fields) != M {
			return nil, fmt.Errorf("Invalid cost board: row %d has %d columns instead of %d", row+1, len(fields), M)
		}

		board.Cost[row] = make([]int, M)
		for col, nums := range fields {
			cost, err := strconv.Atoi(nums)
			if err != nil {
				return nil, fmt.Errorf("Invalid cost value at row %d, col %d: %s", row+1, col+1, nums)
			}
			board.Cost[row][col] = cost
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("Error reading file: %w", err)
	}

	return board, nil
}