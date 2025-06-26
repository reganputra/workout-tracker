package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	"workout-tracker/response"
	"workout-tracker/store"
	"workout-tracker/tokens"
)

type TokenHandler struct {
	tokenStore store.TokenStore
	userStore  store.UserStore
	logger     *log.Logger
}

type createTokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewTokenHandler(tokenStore store.TokenStore, userStore store.UserStore, logger *log.Logger) *TokenHandler {
	return &TokenHandler{
		tokenStore: tokenStore,
		userStore:  userStore,
		logger:     logger,
	}
}

func (th *TokenHandler) HandleCreateToken(w http.ResponseWriter, r *http.Request) {
	var tokenReq createTokenRequest
	err := json.NewDecoder(r.Body).Decode(&tokenReq)
	if err != nil {
		response.BadRequest(w, "Failed to decode token data", err)
	}

	user, err := th.userStore.GetUserByName(tokenReq.Username)
	if err != nil || user == nil {
		response.InternalServerError(w, "Invalid username", err)
		return
	}

	passwordMatch, err := user.PasswordHash.Check(tokenReq.Password)
	if err != nil || !passwordMatch {
		response.InternalServerError(w, "Invalid password", err)
		return
	}

	token, err := th.tokenStore.CreateNewToken(user.Id, 24*time.Hour, tokens.ScopeAuth)
	if err != nil {
		response.InternalServerError(w, "Failed to create token", err)
		return
	}
	response.Success(w, "Token created successfully", map[string]string{"token": token.PlainText})
}
