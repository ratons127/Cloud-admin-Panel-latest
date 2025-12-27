package store

import (
	"crypto/rand"
	"encoding/base64"
	"time"
)

func CreateOAuthState(db *DB) (string, error) {
	if db == nil || db.Conn == nil {
		return "", nil
	}
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", err
	}
	state := base64.RawURLEncoding.EncodeToString(raw)
	hash := HashToken(state)
	expiresAt := time.Now().Add(10 * time.Minute)
	if _, err := db.Conn.Exec(`INSERT INTO oauth_states (token_hash, expires_at) VALUES ($1,$2)`, hash, expiresAt); err != nil {
		return "", err
	}
	return state, nil
}

func ConsumeOAuthState(db *DB, state string) (bool, error) {
	if db == nil || db.Conn == nil {
		return false, nil
	}
	hash := HashToken(state)
	res, err := db.Conn.Exec(`DELETE FROM oauth_states WHERE token_hash = $1 AND expires_at > NOW()`, hash)
	if err != nil {
		return false, err
	}
	if rows, err := res.RowsAffected(); err == nil && rows > 0 {
		return true, nil
	}
	return false, nil
}
