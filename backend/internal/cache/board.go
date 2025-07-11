package cache

import (
	"sync"
)

// Board represents the state of a chess game in memory.
type Board struct {
	Squares  [8][8]string // Each square holds a piece string (e.g., "white-pawn", "black-king", or "")
	LastMove string       // "white" or "black" (whose turn just played)
}

// boardCache is the in-memory map of session string to Board pointer.
var (
	boardCache   = make(map[string]*Board)
	boardCacheMu sync.RWMutex
)

// GetBoard retrieves the board for a session string. Returns nil if not found.
func GetBoard(session string) *Board {
	boardCacheMu.RLock()
	defer boardCacheMu.RUnlock()
	return boardCache[session]
}

// SetBoard sets or updates the board for a session string.
func SetBoard(session string, board *Board) {
	boardCacheMu.Lock()
	defer boardCacheMu.Unlock()
	boardCache[session] = board
}

// NewInitialBoard returns a new Board with the standard chess starting position and last move as "black" (so white moves first).
func NewInitialBoard() *Board {
	var b Board
	// Set up empty squares
	for i := range b.Squares {
		for j := range b.Squares[i] {
			b.Squares[i][j] = ""
		}
	}
	// Place pawns
	for i := 0; i < 8; i++ {
		b.Squares[1][i] = "black-pawn"
		b.Squares[6][i] = "white-pawn"
	}
	// Place other pieces
	b.Squares[0][0], b.Squares[0][7] = "black-rook", "black-rook"
	b.Squares[0][1], b.Squares[0][6] = "black-knight", "black-knight"
	b.Squares[0][2], b.Squares[0][5] = "black-bishop", "black-bishop"
	b.Squares[0][3] = "black-queen"
	b.Squares[0][4] = "black-king"

	b.Squares[7][0], b.Squares[7][7] = "white-rook", "white-rook"
	b.Squares[7][1], b.Squares[7][6] = "white-knight", "white-knight"
	b.Squares[7][2], b.Squares[7][5] = "white-bishop", "white-bishop"
	b.Squares[7][3] = "white-queen"
	b.Squares[7][4] = "white-king"

	b.LastMove = "black" // So white moves first
	return &b
}
