package tokens

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"time"
)

const (
	ScopeAuth = "authentication"
)

type Token struct {
	PlainText string    `json:"plaintext"`
	Hash      []byte    `json:"-"`
	UserID    int       `json:"-"`
	Expired   time.Time `json:"expired"`
	Scopes    string    `json:"-"`
}

func GenerateToken(userId int, ttl time.Duration, scope string) (*Token, error) {
	token := &Token{
		UserID:  userId,
		Expired: time.Now().Add(ttl),
		Scopes:  scope,
	}

	emptyByte := make([]byte, 32)
	_, err := rand.Read(emptyByte)
	if err != nil {
		return nil, err
	}
	token.PlainText = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(emptyByte)
	hash := sha256.Sum256([]byte(token.PlainText))
	token.Hash = hash[:]
	return token, nil
}
