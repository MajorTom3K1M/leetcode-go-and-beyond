package main

func isValidSudokuSolution(board [][]byte) bool {
	rowLen := len(board)
	colLen := len(board[0])

	for i := 0; i < rowLen; i++ {
		for j := 0; j < colLen; j++ {
			for k := 0; k < colLen; k++ {
				if board[i][j] != '.' || board[i][k] != '.' {
					if j != k && board[i][j] == board[i][k] {
						return false
					}
				}
			}
		}
	}

	for i := 0; i < colLen; i++ {
		for j := 0; j < rowLen; j++ {
			for k := 0; k < rowLen; k++ {
				if board[j][i] != '.' || board[k][i] != '.' {
					if j != k && board[j][i] == board[k][i] {
						return false
					}
				}
			}
		}
	}

	for boxRow := 0; boxRow < 3; boxRow++ {
		for boxCol := 0; boxCol < 3; boxCol++ {
			seen := make(map[byte]bool)
			for i := 0; i < 3; i++ {
				for j := 0; j < 3; j++ {
					row := boxRow*3 + i
					col := boxCol*3 + j
					val := board[row][col]
					if val != '.' {
						if seen[val] {
							return false
						}
						seen[val] = true
					}
				}
			}
		}
	}
	return true
}
