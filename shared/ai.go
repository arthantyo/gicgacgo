package shared

import "math"

// AI represents the computer opponent
type AI struct{}

// Minimax implements the minimax algorithm for tic-tac-toe
// Returns the best score for the current player
func (ai *AI) Minimax(board Board, depth int, isMaximizing bool, aiSymbol string, playerSymbol string) int {
	// Check for terminal states
	winner, won := CheckWin(board)
	if won {
		if winner == aiSymbol {
			return 10 - depth // Prefer faster wins
		}
		return depth - 10 // Prefer slower losses
	}

	if CheckDraw(board) {
		return 0
	}

	if isMaximizing {
		bestScore := math.MinInt32
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				if board[i][j] == "" {
					board[i][j] = aiSymbol
					score := ai.Minimax(board, depth+1, false, aiSymbol, playerSymbol)
					board[i][j] = ""
					if score > bestScore {
						bestScore = score
					}
				}
			}
		}
		return bestScore
	} else {
		bestScore := math.MaxInt32
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				if board[i][j] == "" {
					board[i][j] = playerSymbol
					score := ai.Minimax(board, depth+1, true, aiSymbol, playerSymbol)
					board[i][j] = ""
					if score < bestScore {
						bestScore = score
					}
				}
			}
		}
		return bestScore
	}
}

// GetBestMove returns the best move for the AI using minimax
func (ai *AI) GetBestMove(board Board, aiSymbol string, playerSymbol string) (int, int) {
	bestScore := math.MinInt32
	bestRow := -1
	bestCol := -1

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if board[i][j] == "" {
				board[i][j] = aiSymbol
				score := ai.Minimax(board, 0, false, aiSymbol, playerSymbol)
				board[i][j] = ""

				if score > bestScore {
					bestScore = score
					bestRow = i
					bestCol = j
				}
			}
		}
	}

	return bestRow, bestCol
}
