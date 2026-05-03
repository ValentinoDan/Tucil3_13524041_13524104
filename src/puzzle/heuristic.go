package puzzle

func CalculateHeuristic(board *Board, state State) int {
	var target Point
	if state.NextNum < len(board.Checkpoint) {
		target = board.Checkpoint[state.NextNum]
	} else {
		target = board.Goal
	}
	dr := state.Pos.Row - target.Row
	if dr < 0 {
		dr = -dr
	}
	dc := state.Pos.Col - target.Col
	if dc < 0 {
		dc = -dc
	}
	return dr + dc
}
