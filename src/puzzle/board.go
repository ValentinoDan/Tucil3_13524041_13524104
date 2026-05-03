package puzzle

import "fmt"

type Board struct {
	N, M        int
	Grid        [][]rune
	Cost        [][]int
	Start, Goal Point
	Checkpoint  map[int]Point
	Lava        []Point
}

type Point struct {
	Row, Col int
}

// Checker functions
func IsWall(p Point, board *Board) bool {
	return (board.Grid[p.Row][p.Col] == 'X')
}

func IsLava(p Point, board *Board) bool {
	for _, l := range board.Lava {
		if l.Row == p.Row && l.Col == p.Col {
			return true
		}
	}
	return false
}

func IsGoal(p Point, board *Board) bool {
	return (board.Goal.Row == p.Row && board.Goal.Col == p.Col)
}

func IsInBounds(p Point, board *Board) bool {
	return (p.Row >= 0 && p.Row < board.N && p.Col >= 0 && p.Col < board.M)
}

func GetCost(p Point, board *Board) int {
	return board.Cost[p.Row][p.Col]
}

// print board 
func (b *Board) Print() {
	fmt.Printf("Board Size: %d x %d\n", b.N, b.M)
	fmt.Println("Grid:")
	for _, row := range b.Grid {
		for _, char := range row {
			fmt.Printf("%c ", char)
		}
		fmt.Println()
	}
	fmt.Println("Costs:")
	for _, row := range b.Cost {
		for _, cost := range row {
			fmt.Printf("%d ", cost)
		}
		fmt.Println()
	}
}
