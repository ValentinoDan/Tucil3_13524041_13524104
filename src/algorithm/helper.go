package algorithm

import (
	"iceSlidingPuzzle/src/puzzle"
)

// Abs returns the absolute value of x.
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// slide until hits wall / stop
func slide(startState puzzle.State, dir puzzle.Direction, board *puzzle.Board) (puzzle.State, int, bool) {
	currPos := startState.Pos
	nextNum := startState.NextNum
	cost := 0

	dRow, dCol := 0, 0
	switch dir {
	case puzzle.Up:
		dRow = -1
	case puzzle.Down:
		dRow = 1
	case puzzle.Left:
		dCol = -1
	case puzzle.Right:
		dCol = 1
	}

	moved := false

	for {
		nextPos := puzzle.Point{Row: currPos.Row + dRow, Col: currPos.Col + dCol}

		// out of bounds
		if !puzzle.IsInBounds(nextPos, board) {
			return puzzle.State{}, 0, false
		}

		// hits wall
		if puzzle.IsWall(nextPos, board) {
			break
		}

		// continue sliding
		currPos = nextPos
		moved = true
		cost += puzzle.GetCost(currPos, board)

		// lava
		if puzzle.IsLava(currPos, board) {
			return puzzle.State{}, 0, false
		}

		// checkpoint
		for num, pt := range board.Checkpoint {
			if currPos == pt {
				if num == nextNum {
					nextNum++
				} else if num > nextNum { // not sequential
					return puzzle.State{}, 0, false
				}
				break
			}
		}

		// goal
		if puzzle.IsGoal(currPos, board) {
			break
		}
	}

	newState := puzzle.State{
		Pos:     currPos,
		NextNum: nextNum,
	}

	return newState, cost, moved
}
