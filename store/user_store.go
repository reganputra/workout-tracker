package store

import (
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type password struct {
	plainText string
	hash      []byte
}

func (p *password) Set(plainTextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword), 10)
	if err != nil {
		return err
	}
	p.plainText = plainTextPassword
	p.hash = hash
	return nil
}

func (p *password) Check(plainTextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plainTextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

type User struct {
	Id           int       `json:"id"`
	UserName     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash password  `json:"-"`
	Bio          string    `json:"bio"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

var AnonymousUser = &User{}

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

type PostgresUserStore struct {
	db *sql.DB
}

func NewPostgresUserStore(db *sql.DB) *PostgresUserStore {
	return &PostgresUserStore{db: db}
}

type UserStore interface {
	CreateUser(*User) error
	GetUserByName(username string) (*User, error)
	UpdateUser(*User) error
	GetUserToken(scope, tokenPlaintextPassword string) (*User, error)
}

func (store *PostgresUserStore) CreateUser(user *User) error {
	query := "INSERT INTO users (username, email, password_hash, bio) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at"

	err := store.db.QueryRow(query, user.UserName, user.Email, user.PasswordHash.hash, user.Bio).Scan(&user.Id, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (store *PostgresUserStore) GetUserByName(username string) (*User, error) {
	user := &User{
		PasswordHash: password{},
	}
	query := "SELECT id, username, email, password_hash, bio, created_at, updated_at FROM users WHERE username = $1"

	err := store.db.QueryRow(query, username).Scan(&user.Id,
		&user.UserName,
		&user.Email,
		&user.PasswordHash.hash,
		&user.Bio,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (store *PostgresUserStore) UpdateUser(user *User) error {
	query := "UPDATE users SET username = $1, email = $2, bio = $3, updated_at = CURRENT_TIMESTAMP WHERE id = $4 RETURNING updated_at"

	result, err := store.db.Exec(query, user.UserName, user.Email, user.Bio, user.Id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (store *PostgresUserStore) GetUserToken(scope, tokenPlaintextPassword string) (*User, error) {
	tokenHash := sha256.Sum256([]byte(tokenPlaintextPassword))

	fmt.Printf("Looking for token with scope: %s\n", scope)
	fmt.Printf("Token hash: %x\n", tokenHash[:])
	fmt.Printf("Current time: %v\n", time.Now())

	query := "SELECT u.id, u.username, u.email, u.password_hash, u.bio, u.created_at, u.updated_at " +
		"FROM users u INNER JOIN tokens t ON t.user_id = u.id WHERE t.hash = $1 AND t.scope = $2 AND t.expired > $3"

	user := &User{
		PasswordHash: password{},
	}
	err := store.db.QueryRow(query, tokenHash[:], scope, time.Now()).Scan(
		&user.Id,
		&user.UserName,
		&user.Email,
		&user.PasswordHash.hash,
		&user.Bio,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		fmt.Println("No token found in database")
		return nil, nil
	}
	if err != nil {
		fmt.Printf("Database error: %v\n", err)
		return nil, err
	}
	fmt.Printf("Found user: %s\n", user.UserName)
	return user, nil
}
