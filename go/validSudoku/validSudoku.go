package main

import "fmt"

func isValidSudoku(board [][]byte) bool {
	// check row
	for i := range 9 {
		set := make(map[byte]struct{})
		for j := range 9 {
			_, ok := set[board[i][j]]
			if ok {
				return false
			} else if board[i][j] != '.' {
				set[board[i][j]] = struct{}{}
			}
		}
	}

	// check col
	for i := range 9 {
		set := make(map[byte]struct{})
		for j := range 9 {
			_, ok := set[board[j][i]]
			if ok {
				return false
			} else if board[j][i] != '.' {
				set[board[j][i]] = struct{}{}
			}
		}
	}

	// check sub-boxes
	hashSet := make(map[[2]int]map[byte]struct{})
	for i := range 9 {
		for j := range 9 {
			row := i / 3
			col := j / 3
			_, ok := hashSet[[2]int{row, col}][board[i][j]]
			if ok {
				return false
			} else if board[i][j] != '.' {
				_, ok := hashSet[[2]int{row, col}]
				if !ok {
					hashSet[[2]int{row, col}] = make(map[byte]struct{})
				}
				hashSet[[2]int{row, col}][board[i][j]] = struct{}{}
			}
		}
	}

	return true
}

func main() {
	board := [][]byte{
		{'5', '3', '.', '.', '7', '.', '.', '.', '.'},
		{'6', '.', '.', '1', '9', '5', '.', '.', '.'},
		{'.', '9', '8', '.', '.', '.', '.', '6', '.'},
		{'8', '.', '.', '.', '6', '.', '.', '.', '3'},
		{'4', '.', '.', '8', '.', '3', '.', '.', '1'},
		{'7', '.', '.', '.', '2', '.', '.', '.', '6'},
		{'.', '6', '.', '.', '.', '.', '2', '8', '.'},
		{'.', '.', '.', '4', '1', '9', '.', '.', '5'},
		{'.', '.', '.', '.', '8', '.', '.', '7', '9'},
	}
	fmt.Println(isValidSudoku(board))
}
