package puzzle

type Board struct {
	N, M        int
	Grid        [][]rune
	Cost        [][]int
	Start, Goal Point
	Obstacle    map[int]Point
	Lava        []Point
}

type Point struct {
	Row, Col int
}
