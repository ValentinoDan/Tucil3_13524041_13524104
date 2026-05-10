package filehandler

import (
	"fmt"
	"os"
	"strings"

	"iceSlidingPuzzle/src/puzzle"
)

func getDirectionStr(from, to puzzle.Point) string {
	if to.Row < from.Row {
		return "U"
	} else if to.Row > from.Row {
		return "D"
	} else if to.Col < from.Col {
		return "L"
	} else if to.Col > from.Col {
		return "R"
	}
	return "?"
}

func getBoardStateStr(board *puzzle.Board, currentPos puzzle.Point) string {
	var sb strings.Builder
	for r := 0; r < board.N; r++ {
		for c := 0; c < board.M; c++ {
			if r == currentPos.Row && c == currentPos.Col {
				sb.WriteString("Z")
			} else if r == board.Start.Row && c == board.Start.Col {
				sb.WriteString("*")
			} else {
				sb.WriteRune(board.Grid[r][c])
			}
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func SaveOutputTxt(board *puzzle.Board, pathTaken []puzzle.Point, cost int, totalNodes int, executionTime float64, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	var steps []string
	for i := 1; i < len(pathTaken); i++ {
		steps = append(steps, getDirectionStr(pathTaken[i-1], pathTaken[i]))
	}
	file.WriteString(fmt.Sprintf("Solusi gerakan yang ditemukan: %s\n", strings.Join(steps, ", ")))
	file.WriteString(fmt.Sprintf("Cost dari solusi: %d\n", cost))

	if len(pathTaken) > 0 {
		file.WriteString("Initial\n")
		file.WriteString(getBoardStateStr(board, pathTaken[0]))
		file.WriteString("\n")
	}

	for i := 1; i < len(pathTaken); i++ {
		dir := getDirectionStr(pathTaken[i-1], pathTaken[i])
		file.WriteString(fmt.Sprintf("Step %d : %s\n", i, dir))
		file.WriteString(getBoardStateStr(board, pathTaken[i]))
		file.WriteString("\n")
	}
	file.WriteString(fmt.Sprintf("Waktu eksekusi: %.3f ms\n", executionTime))
	file.WriteString(fmt.Sprintf("Banyak iterasi yang dilakukan: %d iterasi", totalNodes))
	return nil
}
