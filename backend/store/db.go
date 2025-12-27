package store

import (
	"crypto/sha256"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

type DB struct {
	Conn *sql.DB
}

func Connect(dbURL string) (*DB, error) {
	if strings.TrimSpace(dbURL) == "" {
		return nil, errors.New("DB_URL is required")
	}
	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}
	conn.SetMaxIdleConns(5)
	conn.SetMaxOpenConns(20)
	conn.SetConnMaxLifetime(30 * time.Minute)
	if err := conn.Ping(); err != nil {
		return nil, err
	}
	return &DB{Conn: conn}, nil
}

func (db *DB) Close() error {
	return db.Conn.Close()
}

func (db *DB) EnsureMigrations(path string) error {
	if _, err := db.Conn.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (version TEXT PRIMARY KEY)`); err != nil {
		return err
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	var migrations []string
	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".sql") {
			continue
		}
		migrations = append(migrations, f.Name())
	}
	sort.Strings(migrations)

	for _, name := range migrations {
		if applied, err := db.isMigrationApplied(name); err != nil {
			return err
		} else if applied {
			continue
		}

		content, err := ioutil.ReadFile(filepath.Join(path, name))
		if err != nil {
			return err
		}
		if _, err := db.Conn.Exec(string(content)); err != nil {
			return err
		}
		if _, err := db.Conn.Exec(`INSERT INTO schema_migrations (version) VALUES ($1)`, name); err != nil {
			return err
		}
	}
	return nil
}

func (db *DB) isMigrationApplied(version string) (bool, error) {
	var v string
	err := db.Conn.QueryRow(`SELECT version FROM schema_migrations WHERE version = $1`, version).Scan(&v)
	if err == sql.ErrNoRows {
		return false, nil
	}
	return err == nil, err
}

func HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func NewInviteToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
