package store

import (
	"database/sql"
	"time"
	"workout-tracker/tokens"
)

type PostgresTokenStore struct {
	db *sql.DB
}

func NewPostgresTokenStore(db *sql.DB) *PostgresTokenStore {
	return &PostgresTokenStore{db: db}
}

type TokenStore interface {
	Insert(token *tokens.Token) error
	CreateNewToken(userId int, ttl time.Duration, scope string) (*tokens.Token, error)
	DeleteAllTokens(userId int, scope string) error
}

func (s *PostgresTokenStore) CreateNewToken(userId int, ttl time.Duration, scope string) (*tokens.Token, error) {
	token, err := tokens.GenerateToken(userId, ttl, scope)
	if err != nil {
		return nil, err
	}
	err = s.Insert(token)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (s *PostgresTokenStore) Insert(token *tokens.Token) error {
	query := `INSERT INTO tokens (hash, user_id, expired, scope) VALUES ($1, $2, $3, $4)`
	_, err := s.db.Exec(query, token.Hash, token.UserID, token.Expired, token.Scopes)
	return err
}

func (s *PostgresTokenStore) DeleteAllTokens(userId int, scope string) error {
	query := `DELETE FROM tokens WHERE scope = $1 AND user_id = $2`
	_, err := s.db.Exec(query, scope, userId)
	return err
}
