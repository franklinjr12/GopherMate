package movevalidation

import (
	"errors"
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
func ValidateMove(move MoveData) (bool, error) {
	switch {
	case isPiece(move.Piece, "pawn"):
		return validatePawnMove(move)
	case isPiece(move.Piece, "rook"):
		return validateRookMove(move)
	case isPiece(move.Piece, "knight"):
		return validateKnightMove(move)
	case isPiece(move.Piece, "bishop"):
		return validateBishopMove(move)
	case isPiece(move.Piece, "queen"):
		return validateQueenMove(move)
	case isPiece(move.Piece, "king"):
		return validateKingMove(move)
	default:
		return false, errors.New("Unknown piece type")
	}
}

func isPiece(piece, name string) bool {
	return piece == name || piece == "white-"+name || piece == "black-"+name
}

// validateKingMove validates king moves (one square in any direction, no castling yet)
func validateKingMove(move MoveData) (bool, error) {
	deltaRow := abs(move.To.Row - move.From.Row)
	deltaCol := abs(move.To.Col - move.From.Col)
	if (deltaRow <= 1 && deltaCol <= 1) && (deltaRow != 0 || deltaCol != 0) {
		return true, nil
	}
	return false, errors.New("Invalid king move")
}

// validateQueenMove validates queen moves (combines rook and bishop logic)
func validateQueenMove(move MoveData) (bool, error) {
	// Queen moves like rook or bishop
	okRook, _ := validateRookMove(move)
	okBishop, _ := validateBishopMove(move)
	if okRook || okBishop {
		return true, nil
	}
	return false, errors.New("Invalid queen move")
}

// validateBishopMove validates bishop moves (diagonal, any distance, no jumping check)
func validateBishopMove(move MoveData) (bool, error) {
	deltaRow := abs(move.To.Row - move.From.Row)
	deltaCol := abs(move.To.Col - move.From.Col)
	// Bishop moves only if row and col change by the same amount (diagonal)
	if deltaRow == deltaCol && deltaRow != 0 {
		return true, nil
	}
	return false, errors.New("Invalid bishop move")
}

// validateKnightMove validates knight (horse) moves (L-shape: 2 by 1 or 1 by 2)
func validateKnightMove(move MoveData) (bool, error) {
	deltaRow := abs(move.To.Row - move.From.Row)
	deltaCol := abs(move.To.Col - move.From.Col)
	if (deltaRow == 2 && deltaCol == 1) || (deltaRow == 1 && deltaCol == 2) {
		return true, nil
	}
	return false, errors.New("Invalid knight move")
}

// validateRookMove validates rook moves (horizontal or vertical, any distance, no jumping check)
func validateRookMove(move MoveData) (bool, error) {
	deltaRow := move.To.Row - move.From.Row
	deltaCol := move.To.Col - move.From.Col
	// Rook moves only if either row or col is unchanged, but not both
	if (deltaRow == 0 && deltaCol != 0) || (deltaCol == 0 && deltaRow != 0) {
		return true, nil
	}
	return false, errors.New("Invalid rook move")
}

// validatePawnMove validates pawn moves (basic forward, capture, double move, no en passant yet)
func validatePawnMove(move MoveData) (bool, error) {
	// White pawns move up (row decreases), black pawns move down (row increases)
	rowDir := 1
	startRow := 1
	if move.Piece == "white-pawn" {
		rowDir = -1
		startRow = 6
	}
	deltaRow := move.To.Row - move.From.Row
	deltaCol := move.To.Col - move.From.Col

	// Forward move
	if deltaCol == 0 {
		if deltaRow == rowDir {
			return true, nil // single forward
		}
		if move.From.Row == startRow && deltaRow == 2*rowDir {
			return true, nil // double forward from start
		}
	}
	// Capture (diagonal)
	if abs(deltaCol) == 1 && deltaRow == rowDir {
		// In real chess, must check if target square has opponent piece
		return true, nil // allow for now
	}
	return false, errors.New("Invalid pawn move")
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
