package algorithm

import (
	"iceSlidingPuzzle/src/puzzle"
)

func UCS(board *puzzle.Board) (*puzzle.Node, int){
	pq := &puzzle.PriorityQueue{}
	visited := make(map[puzzle.State]bool)	
	iter := 0

	startNode := &puzzle.Node{
		State: puzzle.State{Pos: board.Start, NextNum: 0},
		Cost:  0,
		Depth: 0,
		Parent: nil,
		Dir: puzzle.Nil,
	}
	pq.Push(startNode)

	for !pq.IsEmpty() {
		curr := pq.Pop()
		iter++

		if visited[curr.State]{
			continue
		}
		visited[curr.State] = true

		// check if its target found already & all checkpoint have been passed
		if puzzle.IsGoal(curr.State.Pos, board) && curr.State.NextNum == len(board.Checkpoint) {
			return curr, iter
		}

		// explore others
		for _, dir := range []puzzle.Direction{puzzle.Up, puzzle.Down, puzzle.Left, puzzle.Right} {
			nextState, moveCost, valid := slide(curr.State, dir, board)
			
			if valid {
				nextNode := &puzzle.Node{
					State:  nextState,
					Cost:   curr.Cost + moveCost,
					Depth:  curr.Depth + 1,
					Parent: curr,
					Dir:    dir,
				}
				pq.Push(nextNode)
			}
		}
	}
	return nil, iter
}