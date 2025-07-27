package cache

import (
	"sync"
	"time"
)

// Board represents the state of a chess game in memory.
type Board struct {
	Squares          [8][8]string // Each square holds a piece string (e.g., "white-pawn", "black-king", or "")
	LastMove         string       // "white" or "black" (whose turn just played)
	LastMoveNumber   int          // The move number of the last move made
	LastMoveNotation string       // The notation of the last move made
	DrawOffer        string       // "white", "black", or "" (who offered draw, empty if none)
	DrawOfferPending bool         // true if a draw offer is pending, false otherwise
}

// boardCache is the in-memory map of session string to Board pointer and its last updated time.
var (
	boardCache   = make(map[string]*boardCacheEntry)
	boardCacheMu sync.RWMutex
)

type boardCacheEntry struct {
	Board     *Board
	UpdatedAt time.Time
}

// GetBoard retrieves the board for a session string. Returns nil if not found or expired.
func GetBoard(session string) *Board {
	boardCacheMu.RLock()
	entry := boardCache[session]
	boardCacheMu.RUnlock()
	if entry == nil {
		return nil
	}
	if time.Since(entry.UpdatedAt) > 30*time.Minute {
		return nil
	}
	return entry.Board
}

// SetBoard sets or updates the board for a session string.
func SetBoard(session string, board *Board) {
	boardCacheMu.Lock()
	boardCache[session] = &boardCacheEntry{
		Board:     board,
		UpdatedAt: time.Now(),
	}
	boardCacheMu.Unlock()
}

// CleanExpiredBoards removes boards that have not been updated in the last 30 minutes.
func CleanExpiredBoards() {
	boardCacheMu.Lock()
	now := time.Now()
	for session, entry := range boardCache {
		if now.Sub(entry.UpdatedAt) > 30*time.Minute {
			delete(boardCache, session)
		}
	}
	boardCacheMu.Unlock()
}

// ClearBoard removes a specific board from the cache.
func ClearBoard(session string) {
	boardCacheMu.Lock()
	delete(boardCache, session)
	boardCacheMu.Unlock()
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

	b.LastMove = "black"    // So white moves first
	b.LastMoveNumber = 0    // No moves made yet
	b.LastMoveNotation = "" // No moves made yet
	return &b
}
