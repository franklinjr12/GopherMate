package movevalidation

import (
	"errors"
	"gophermatebackend/internal/cache"
)

type Position struct {
	Row int
	Col int
}

type MoveData struct {
	Piece string
	From  Position
	To    Position
}

// ValidateMove is the entrypoint for move validation. It dispatches to the correct piece validator.
func ValidateMove(board *cache.Board, move MoveData) (bool, error) {
	// check if the move is from the opposing player
	if (board.LastMove == "white" && move.Piece[:5] == "white") ||
		(board.LastMove == "black" && move.Piece[:5] == "black") {
		return false, errors.New("it's not your turn")
	}
	switch {
	case isPiece(move.Piece, "pawn"):
		return validatePawnMove(board, move)
	case isPiece(move.Piece, "rook"):
		return validateRookMove(board, move)
	case isPiece(move.Piece, "knight"):
		return validateKnightMove(board, move)
	case isPiece(move.Piece, "bishop"):
		return validateBishopMove(board, move)
	case isPiece(move.Piece, "queen"):
		return validateQueenMove(board, move)
	case isPiece(move.Piece, "king"):
		return validateKingMove(board, move)
	default:
		return false, errors.New("Unknown piece type")
	}
}

func isPiece(piece, name string) bool {
	return piece == name || piece == "white-"+name || piece == "black-"+name
}

// validateKingMove validates king moves (one square in any direction, no castling yet)
func validateKingMove(board *cache.Board, move MoveData) (bool, error) {
	deltaRow := abs(move.To.Row - move.From.Row)
	deltaCol := abs(move.To.Col - move.From.Col)
	if (deltaRow <= 1 && deltaCol <= 1) && (deltaRow != 0 || deltaCol != 0) {
		dest := board.Squares[move.To.Row][move.To.Col]
		if dest == "" || isOpponentPiece(dest, getColor(move.Piece)) {
			return true, nil
		}
		return false, errors.New("King cannot move to a square occupied by friendly piece")
	}
	return false, errors.New("Invalid king move")
}

// validateQueenMove validates queen moves (combines rook and bishop logic)
func validateQueenMove(board *cache.Board, move MoveData) (bool, error) {
	// Queen moves like rook or bishop
	okRook, _ := validateRookMove(board, move)
	if okRook {
		return true, nil
	}
	okBishop, _ := validateBishopMove(board, move)
	if okBishop {
		return true, nil
	}
	return false, errors.New("Invalid queen move")
}

// validateBishopMove validates bishop moves (diagonal, any distance, no jumping check)
func validateBishopMove(board *cache.Board, move MoveData) (bool, error) {
	deltaRow := move.To.Row - move.From.Row
	deltaCol := move.To.Col - move.From.Col
	if abs(deltaRow) == abs(deltaCol) && deltaRow != 0 {
		rowStep := 1
		if deltaRow < 0 {
			rowStep = -1
		}
		colStep := 1
		if deltaCol < 0 {
			colStep = -1
		}
		r, c := move.From.Row+rowStep, move.From.Col+colStep
		for r != move.To.Row && c != move.To.Col {
			if board.Squares[r][c] != "" {
				return false, errors.New("Bishop cannot jump over pieces")
			}
			r += rowStep
			c += colStep
		}
		dest := board.Squares[move.To.Row][move.To.Col]
		if dest == "" || isOpponentPiece(dest, getColor(move.Piece)) {
			return true, nil
		}
		return false, errors.New("Bishop cannot move to a square occupied by friendly piece")
	}
	return false, errors.New("Invalid bishop move")
}

// validateKnightMove validates knight (horse) moves (L-shape: 2 by 1 or 1 by 2)
func validateKnightMove(board *cache.Board, move MoveData) (bool, error) {
	deltaRow := abs(move.To.Row - move.From.Row)
	deltaCol := abs(move.To.Col - move.From.Col)
	if (deltaRow == 2 && deltaCol == 1) || (deltaRow == 1 && deltaCol == 2) {
		dest := board.Squares[move.To.Row][move.To.Col]
		if dest == "" || isOpponentPiece(dest, getColor(move.Piece)) {
			return true, nil
		}
		return false, errors.New("Knight cannot move to a square occupied by friendly piece")
	}
	return false, errors.New("Invalid knight move")
}

// validateRookMove validates rook moves (horizontal or vertical, any distance, no jumping check)
func validateRookMove(board *cache.Board, move MoveData) (bool, error) {
	deltaRow := move.To.Row - move.From.Row
	deltaCol := move.To.Col - move.From.Col
	if (deltaRow == 0 && deltaCol != 0) || (deltaCol == 0 && deltaRow != 0) {
		rowStep, colStep := 0, 0
		if deltaRow != 0 {
			rowStep = deltaRow / abs(deltaRow)
		}
		if deltaCol != 0 {
			colStep = deltaCol / abs(deltaCol)
		}
		r, c := move.From.Row+rowStep, move.From.Col+colStep
		for r != move.To.Row || c != move.To.Col {
			if board.Squares[r][c] != "" {
				return false, errors.New("Rook cannot jump over pieces")
			}
			r += rowStep
			c += colStep
		}
		dest := board.Squares[move.To.Row][move.To.Col]
		if dest == "" || isOpponentPiece(dest, getColor(move.Piece)) {
			return true, nil
		}
		return false, errors.New("Rook cannot move to a square occupied by friendly piece")
	}
	return false, errors.New("Invalid rook move")
}

// getColor returns the color of the piece ("white" or "black")
func getColor(piece string) string {
	if len(piece) >= 6 && piece[:5] == "white" {
		return "white"
	}
	if len(piece) >= 6 && piece[:5] == "black" {
		return "black"
	}
	return ""
}

// validatePawnMove validates pawn moves (basic forward, capture, double move, no en passant yet)
func validatePawnMove(board *cache.Board, move MoveData) (bool, error) {
	// White pawns move up (row decreases), black pawns move down (row increases)
	rowDir := 1
	startRow := 1
	myColor := "black"
	if move.Piece == "white-pawn" {
		rowDir = -1
		startRow = 6
		myColor = "white"
	}
	deltaRow := move.To.Row - move.From.Row
	deltaCol := move.To.Col - move.From.Col

	// Forward move (no capture)
	if deltaCol == 0 {
		// Single forward
		if deltaRow == rowDir {
			if board.Squares[move.To.Row][move.To.Col] == "" {
				return true, nil
			}
			return false, errors.New("Pawn cannot move forward to occupied square")
		}
		// Double forward from start
		if move.From.Row == startRow && deltaRow == 2*rowDir {
			midRow := move.From.Row + rowDir
			if board.Squares[midRow][move.From.Col] == "" && board.Squares[move.To.Row][move.To.Col] == "" {
				return true, nil
			}
			return false, errors.New("Pawn cannot jump over or land on occupied square")
		}
	}

	// Capture (diagonal)
	if abs(deltaCol) == 1 && deltaRow == rowDir {
		target := board.Squares[move.To.Row][move.To.Col]
		if target != "" && isOpponentPiece(target, myColor) {
			return true, nil
		}
		return false, errors.New("Pawn capture must target opponent piece")
	}

	return false, errors.New("Invalid pawn move")
}

// isOpponentPiece checks if a piece string belongs to the opponent color.
func isOpponentPiece(piece, myColor string) bool {
	if myColor == "white" {
		return len(piece) > 6 && piece[:5] == "black"
	}
	if myColor == "black" {
		return len(piece) > 5 && piece[:5] == "white"
	}
	return false
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
