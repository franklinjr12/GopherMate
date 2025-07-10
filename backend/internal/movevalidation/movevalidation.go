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
	case move.Piece == "white-pawn" || move.Piece == "black-pawn":
		return validatePawnMove(move)
	// Add other pieces here
	default:
		return false, errors.New("Unknown piece type")
	}
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
