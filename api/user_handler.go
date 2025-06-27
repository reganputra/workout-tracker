package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"
	"workout-tracker/response"
	"workout-tracker/store"
)

type registerUserRequest struct {
	UserName string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Bio      string `json:"bio"`
}

type UserHandler struct {
	userStore store.UserStore
	logger    *log.Logger
}

func NewUserHandler(userStore store.UserStore, logger *log.Logger) *UserHandler {
	return &UserHandler{
		userStore: userStore,
		logger:    logger,
	}
}

func (uh *UserHandler) validateUserRequest(req *registerUserRequest) error {
	if req.UserName == "" {
		return errors.New("Username is required")
	}

	if len(req.Password) < 8 {
		return errors.New("Password must be at least 8 characters long")
	}

	if req.Password == "" {
		return errors.New("Password is required")
	}

	if req.Email == "" {
		return errors.New("Email is required")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		return errors.New("Invalid email format")
	}
	return nil
}

func (uh *UserHandler) HandleRegisterUser(w http.ResponseWriter, r *http.Request) {
	var userReq registerUserRequest
	err := json.NewDecoder(r.Body).Decode(&userReq)
	if err != nil {
		response.BadRequest(w, "Failed to decode user data", err)
	}

	err = uh.validateUserRequest(&userReq)
	if err != nil {
		response.BadRequest(w, "Invalid user data", err)
	}
	user := &store.User{
		UserName: userReq.UserName,
		Email:    userReq.Email,
		Bio:      userReq.Bio,
	}

	err = user.PasswordHash.Set(userReq.Password)
	if err != nil {
		response.InternalServerError(w, "Failed hashing password", err)
		return
	}
	err = uh.userStore.CreateUser(user)
	if err != nil {
		response.InternalServerError(w, "Failed to create user", err)
		return
	}
	response.UserCreated(w, user)
}
