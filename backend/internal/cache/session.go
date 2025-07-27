package cache

import (
	"sync"
	"time"
)

// sessionCacheEntry holds userID and expiry for a session token.
type sessionCacheEntry struct {
	UserID    int64
	ExpiresAt time.Time
}

var (
	sessionCache   = make(map[string]*sessionCacheEntry)
	sessionCacheMu sync.RWMutex
)

// GetUserIDByToken returns userID if token is present and not expired.
func GetUserIDByToken(token string) (int64, bool) {
	sessionCacheMu.RLock()
	entry, ok := sessionCache[token]
	sessionCacheMu.RUnlock()
	if !ok || time.Now().After(entry.ExpiresAt) {
		return 0, false
	}
	return entry.UserID, true
}

// SetUserIDForToken sets the userID and expiry for a session token.
func SetUserIDForToken(token string, userID int64, expiresAt time.Time) {
	sessionCacheMu.Lock()
	sessionCache[token] = &sessionCacheEntry{
		UserID:    userID,
		ExpiresAt: expiresAt,
	}
	sessionCacheMu.Unlock()
}

// DeleteUserIDForToken removes a session token from the cache.
func DeleteUserIDForToken(token string) {
	sessionCacheMu.Lock()
	delete(sessionCache, token)
	sessionCacheMu.Unlock()
}

// CleanExpiredSessions removes expired session tokens from the cache.
func CleanExpiredSessions() {
	sessionCacheMu.Lock()
	now := time.Now()
	for token, entry := range sessionCache {
		if now.After(entry.ExpiresAt) {
			delete(sessionCache, token)
		}
	}
	sessionCacheMu.Unlock()
}
