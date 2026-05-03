package puzzle

type Direction string

const (
	Up    Direction = "UP"
	Down  Direction = "DOWN"
	Left  Direction = "LEFT"
	Right Direction = "RIGHT"
	Nil Direction   = "NIL"
)

type Node struct {
	State 	State 
	Cost  	int   
	TurnCount int
	Parent *Node
	Dir 	Direction
}

type State struct {
	Pos 	Point
	NextNum int
}

type PriorityQueue struct {
	Nodes []*Node
}

// add new nodes and make an adjustment
func (pq *PriorityQueue) Push(n *Node) {
	pq.Nodes = append(pq.Nodes, n)
	pq.upHeap(len(pq.Nodes) - 1)
}

// take and remove nodes with min cost
func (pq *PriorityQueue) Pop() *Node {
	if len(pq.Nodes) == 0 {
		return nil
	}

	min := pq.Nodes[0]

	// put last element to the root
	lastIdx := len(pq.Nodes) - 1
	pq.Nodes[0] = pq.Nodes[lastIdx]
	pq.Nodes = pq.Nodes[:lastIdx] 

	// adjust the heap 
	if len(pq.Nodes) > 0 {
		pq.downHeap(0)
	}

	return min
}

// check empty
func (pq *PriorityQueue) IsEmpty() bool {
	return len(pq.Nodes) == 0
}

// pick up node with smallest cost and adjust the heap
func (pq *PriorityQueue) upHeap(index int) {
	for index > 0 {
		parent := (index - 1) / 2

		// child >= parent
		if pq.Nodes[index].Cost >= pq.Nodes[parent].Cost {
			break
		}

		// switch pos with parent
		pq.Nodes[index], pq.Nodes[parent] = pq.Nodes[parent], pq.Nodes[index]
		index = parent
	}
}

// put down node with larger cost and adjust the heap
func (pq *PriorityQueue) downHeap(index int) {
	lastIdx := len(pq.Nodes) - 1
	for {
		left := 2*index + 1
		right := 2*index + 2
		smallest := index

		// check left child
		if left <= lastIdx && pq.Nodes[left].Cost < pq.Nodes[smallest].Cost {
			smallest = left
		}

		// check right child
		if right <= lastIdx && pq.Nodes[right].Cost < pq.Nodes[smallest].Cost {
			smallest = right
		}

		if smallest == index {
			break
		}

		// switch pos with smallest
		pq.Nodes[index], pq.Nodes[smallest] = pq.Nodes[smallest], pq.Nodes[index]
		index = smallest
	}
}
