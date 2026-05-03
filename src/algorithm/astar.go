package algorithm

import (
	"iceSlidingPuzzle/src/puzzle"
)

func AStarSearch(board *puzzle.Board) ([]puzzle.Point, int) {
	pq := &puzzle.PriorityQueue{}

	// gScore: Map untuk menyimpan biaya nyata terendah dari Start ke State tertentu
	// Ini adalah g(n)
	gScore := make(map[puzzle.State]int)

	// visited: Untuk menandai state yang sudah diproses secara optimal
	visited := make(map[puzzle.State]bool)

	startState := puzzle.State{
		Pos:     board.Start,
		NextNum: 0,
	}

	// Inisialisasi awal
	gScore[startState] = 0
	h := calculateHeuristic(board, startState)

	// Node awal dimasukkan ke PQ
	// f(n) = g(n) + h(n) => 0 + h
	pq.Push(&puzzle.Node{
		State:  startState,
		Cost:   0 + h,
		Parent: nil,
	})

	var finalNode *puzzle.Node
	directions := []puzzle.Direction{puzzle.Up, puzzle.Down, puzzle.Left, puzzle.Right}

	for !pq.IsEmpty() {
		// Ambil node dengan f(n) terkecil
		curr := pq.Pop()

		// Cek apakah sudah sampai Goal dengan semua checkpoint
		if puzzle.IsGoal(curr.State.Pos, board) && curr.State.NextNum == len(board.Checkpoint) {
			finalNode = curr
			break
		}

		if visited[curr.State] {
			continue
		}
		visited[curr.State] = true

		for _, dir := range directions {
			// Fungsi slide dari helper.go
			nextState, moveCost, moved := slide(curr.State, dir, board)

			if moved {
				// g(n) baru = g(n) sekarang + biaya langkah barusan
				tentativeGScore := gScore[curr.State] + moveCost

				// Jika rute ini lebih murah dari rute yang pernah dicatat sebelumnya
				if oldG, exists := gScore[nextState]; !exists || tentativeGScore < oldG {
					gScore[nextState] = tentativeGScore

					// h(n) dihitung menggunakan fungsi yang sama dengan GBFS
					h := calculateHeuristic(board, nextState)

					// Push ke PQ dengan Prioritas f(n) = g(n) + h(n)
					pq.Push(&puzzle.Node{
						State:  nextState,
						Cost:   tentativeGScore + h, // KUNCI UTAMA A*
						Parent: curr,
						Dir:    dir,
					})
				}
			}
		}
	}

	if finalNode == nil {
		return []puzzle.Point{}, 0
	}

	// Traceback (Mundur dari Goal ke Start)
	var pathTaken []puzzle.Point
	for n := finalNode; n != nil; n = n.Parent {
		pathTaken = append(pathTaken, n.State.Pos)
	}

	// Membalikkan array (Reverse) karena append menambahkan ke belakang
	for i, j := 0, len(pathTaken)-1; i < j; i, j = i+1, j-1 {
		pathTaken[i], pathTaken[j] = pathTaken[j], pathTaken[i]
	}

	return pathTaken, gScore[finalNode.State]
}
