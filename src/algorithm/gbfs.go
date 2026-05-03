package algorithm

import (
	"iceSlidingPuzzle/src/puzzle"
)

func GreedyBestFirstSearch(board *puzzle.Board) ([]puzzle.Point, int) {
	var pathTaken []puzzle.Point
	var totalCost int

	pq := &puzzle.PriorityQueue{}
	visited := make(map[puzzle.State]bool)
	stateCost := make(map[puzzle.State]int)
	startState := puzzle.State{
		Pos:     board.Start,
		NextNum: 0,
	}
	stateCost[startState] = 0

	startNode := &puzzle.Node{
		State:  startState,
		Cost:   calculateHeuristic(board, startState),
		Depth:  0,
		Parent: nil,
		Dir:    puzzle.Nil,
	}

	pq.Push(startNode)
	visited[startState] = true

	var finalNode *puzzle.Node
	directions := []puzzle.Direction{puzzle.Up, puzzle.Down, puzzle.Left, puzzle.Right}

	for !pq.IsEmpty() {
		curr := pq.Pop()
		if puzzle.IsGoal(curr.State.Pos, board) && curr.State.NextNum == len(board.Checkpoint) {
			finalNode = curr
			break
		}

		for _, dir := range directions {
			nextState, moveCost, moved := slide(curr.State, dir, board)
			if moved && !visited[nextState] {
				visited[nextState] = true
				stateCost[nextState] = stateCost[curr.State] + moveCost

				neighborNode := &puzzle.Node{
					State:  nextState,
					Cost:   calculateHeuristic(board, nextState),
					Depth:  curr.Depth + 1,
					Parent: curr,
					Dir:    dir,
				}

				pq.Push(neighborNode)
			}
		}
	}

	if finalNode == nil {
		return pathTaken, 0
	}
	currNode := finalNode
	for currNode != nil {
		pathTaken = append(pathTaken, currNode.State.Pos)
		currNode = currNode.Parent
	}
	for i, j := 0, len(pathTaken)-1; i < j; i, j = i+1, j-1 {
		pathTaken[i], pathTaken[j] = pathTaken[j], pathTaken[i]
	}

	totalCost = stateCost[finalNode.State]

	return pathTaken, totalCost
}
