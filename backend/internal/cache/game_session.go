package cache

import (
	"fmt"
	"sync"
	"time"
)

// gameSessionCacheEntry holds validation result and expiry for a game-user combination.
type gameSessionCacheEntry struct {
	IsValid   bool
	ExpiresAt time.Time
}

var (
	gameSessionCache   = make(map[string]*gameSessionCacheEntry)
	gameSessionCacheMu sync.RWMutex
)

// GetGameSessionValidation returns the validation result if present and not expired.
func GetGameSessionValidation(key string) (bool, bool) {
	gameSessionCacheMu.RLock()
	entry, ok := gameSessionCache[key]
	gameSessionCacheMu.RUnlock()
	if !ok || time.Now().After(entry.ExpiresAt) {
		return false, false
	}
	return entry.IsValid, true
}

// SetGameSessionValidation sets the validation result and expiry for a game-user combination.
func SetGameSessionValidation(key string, isValid bool, expiresAt time.Time) {
	gameSessionCacheMu.Lock()
	gameSessionCache[key] = &gameSessionCacheEntry{
		IsValid:   isValid,
		ExpiresAt: expiresAt,
	}
	gameSessionCacheMu.Unlock()
}

// DeleteGameSessionValidation removes a game-user combination from the cache.
func DeleteGameSessionValidation(key string) {
	gameSessionCacheMu.Lock()
	delete(gameSessionCache, key)
	gameSessionCacheMu.Unlock()
}

// CleanExpiredGameSessions removes expired game session validations from the cache.
func CleanExpiredGameSessions() {
	gameSessionCacheMu.Lock()
	now := time.Now()
	for key, entry := range gameSessionCache {
		if now.After(entry.ExpiresAt) {
			delete(gameSessionCache, key)
		}
	}
	gameSessionCacheMu.Unlock()
}

// GenerateGameSessionKey creates a unique cache key for a game-user combination.
func GenerateGameSessionKey(gameID string, userID int64) string {
	return fmt.Sprintf("%s:%d", gameID, userID)
}
