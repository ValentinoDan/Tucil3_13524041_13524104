package algorithm

import (
	"iceSlidingPuzzle/src/puzzle"
	"math"
)

const FOUND = -1

func idaStarDFS(board *puzzle.Board, currNode *puzzle.Node, g int, threshold int, pathMap map[puzzle.State]bool) (int, *puzzle.Node) {
	h := puzzle.CalculateHeuristic(board, currNode.State)
	f := g + h
	if f > threshold {
		return f, nil
	}

	if puzzle.IsGoal(currNode.State.Pos, board) && currNode.State.NextNum == len(board.Checkpoint) {
		return FOUND, currNode
	}

	min := math.MaxInt32
	directions := []puzzle.Direction{puzzle.Up, puzzle.Down, puzzle.Left, puzzle.Right}

	for _, dir := range directions {
		nextState, moveCost, moved := slide(currNode.State, dir, board)
		if moved && !pathMap[nextState] {
			pathMap[nextState] = true
			newG := g + moveCost

			neighborNode := &puzzle.Node{
				State:  nextState,
				Cost:   newG + puzzle.CalculateHeuristic(board, nextState),
				Depth:  currNode.Depth + 1,
				Parent: currNode,
				Dir:    dir,
			}

			t, resultNode := idaStarDFS(board, neighborNode, newG, threshold, pathMap)
			if t == FOUND {
				return FOUND, resultNode
			}
			if t < min {
				min = t
			}

			delete(pathMap, nextState)
		}
	}

	return min, nil
}

func IDAStarSearch(board *puzzle.Board) ([]puzzle.Point, int) {
	startState := puzzle.State{
		Pos:     board.Start,
		NextNum: 0,
	}

	startNode := &puzzle.Node{
		State:  startState,
		Cost:   puzzle.CalculateHeuristic(board, startState),
		Depth:  0,
		Parent: nil,
		Dir:    puzzle.Nil,
	}

	threshold := puzzle.CalculateHeuristic(board, startState)

	for {
		pathMap := make(map[puzzle.State]bool)
		pathMap[startState] = true
		t, finalNode := idaStarDFS(board, startNode, 0, threshold, pathMap)

		if t == FOUND {
			var pathTaken []puzzle.Point
			curr := finalNode
			for curr != nil {
				pathTaken = append(pathTaken, curr.State.Pos)
				curr = curr.Parent
			}
			for i, j := 0, len(pathTaken)-1; i < j; i, j = i+1, j-1 {
				pathTaken[i], pathTaken[j] = pathTaken[j], pathTaken[i]
			}
			return pathTaken, finalNode.Cost
		}
		if t == math.MaxInt32 {
			return []puzzle.Point{}, 0
		}
		threshold = t
	}
}
