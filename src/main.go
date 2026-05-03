package main

import (
	"fmt"
	"iceSlidingPuzzle/src/filehandler"
)

func main() {
	board, err := filehandler.ParseBoard("test/test2.txt")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	board.Print()
}
