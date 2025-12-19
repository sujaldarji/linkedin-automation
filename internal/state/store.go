package state

import (
	"database/sql"
	"errors"
	"time"

	_ "modernc.org/sqlite"

)

type Store struct {
	db *sql.DB
}

// NewStore opens (or creates) the SQLite database and ensures schema
func NewStore(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	if err := migrate(db); err != nil {
		return nil, err
	}

	return &Store{db: db}, nil
}

func migrate(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS profiles (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		profile_url TEXT UNIQUE NOT NULL,
		visited BOOLEAN NOT NULL DEFAULT 0,
		connection_sent BOOLEAN NOT NULL DEFAULT 0,
		status TEXT NOT NULL DEFAULT 'new',
		created_at DATETIME NOT NULL,
		last_action_at DATETIME
	);
	`
	_, err := db.Exec(schema)
	return err
}

// EnsureProfile inserts a profile if it does not already exist
func (s *Store) EnsureProfile(profileURL string) error {
	_, err := s.db.Exec(
		`INSERT OR IGNORE INTO profiles (profile_url, created_at) VALUES (?, ?)`,
		profileURL,
		time.Now(),
	)
	return err
}

// CanVisit returns true if profile has not been visited or failed
func (s *Store) CanVisit(profileURL string) (bool, error) {
	row := s.db.QueryRow(
		`SELECT visited, status FROM profiles WHERE profile_url = ?`,
		profileURL,
	)

	var visited bool
	var status string

	if err := row.Scan(&visited, &status); err != nil {
		return false, err
	}

	if status == "failed" {
		return false, nil
	}

	return !visited, nil
}

// CanSendConnection checks if a connection request is allowed
func (s *Store) CanSendConnection(profileURL string) (bool, error) {
	row := s.db.QueryRow(
		`SELECT visited, connection_sent, status FROM profiles WHERE profile_url = ?`,
		profileURL,
	)

	var visited, sent bool
	var status string

	if err := row.Scan(&visited, &sent, &status); err != nil {
		return false, err
	}

	if !visited || sent {
		return false, nil
	}

	if status == "skipped" || status == "failed" {
		return false, nil
	}

	return true, nil
}

// MarkVisited marks a profile as visited
func (s *Store) MarkVisited(profileURL string) error {
	now := time.Now()
	_, err := s.db.Exec(
		`UPDATE profiles
		 SET visited = 1, status = 'visited', last_action_at = ?
		 WHERE profile_url = ?`,
		now,
		profileURL,
	)
	return err
}

// MarkConnectionSent marks a connection request as sent
func (s *Store) MarkConnectionSent(profileURL string) error {
	now := time.Now()
	_, err := s.db.Exec(
		`UPDATE profiles
		 SET connection_sent = 1, status = 'connected', last_action_at = ?
		 WHERE profile_url = ?`,
		now,
		profileURL,
	)
	return err
}

// MarkSkipped marks a profile as skipped
func (s *Store) MarkSkipped(profileURL string) error {
	now := time.Now()
	_, err := s.db.Exec(
		`UPDATE profiles
		 SET status = 'skipped', last_action_at = ?
		 WHERE profile_url = ?`,
		now,
		profileURL,
	)
	return err
}

// MarkFailed marks a profile as failed and prevents future attempts
func (s *Store) MarkFailed(profileURL string) error {
	now := time.Now()
	_, err := s.db.Exec(
		`UPDATE profiles
		 SET status = 'failed', last_action_at = ?
		 WHERE profile_url = ?`,
		now,
		profileURL,
	)
	return err
}

// GetPendingProfiles returns profiles that have not been visited yet
func (s *Store) GetPendingProfiles(limit int) ([]ProfileState, error) {
	rows, err := s.db.Query(
		`SELECT profile_url, visited, connection_sent, status, created_at, last_action_at
		 FROM profiles
		 WHERE visited = 0 AND status = 'new'
		 ORDER BY created_at
		 LIMIT ?`,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []ProfileState
	for rows.Next() {
		var p ProfileState
		if err := rows.Scan(
			&p.ProfileURL,
			&p.Visited,
			&p.ConnectionSent,
			&p.Status,
			&p.CreatedAt,
			&p.LastActionAt,
		); err != nil {
			return nil, err
		}
		results = append(results, p)
	}

	return results, nil
}

// Close closes the underlying database
func (s *Store) Close() error {
	if s.db == nil {
		return errors.New("store not initialized")
	}
	return s.db.Close()
}
