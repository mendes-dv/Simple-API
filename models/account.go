package models

import (
	"golang.org/x/crypto/bcrypt"
	"math/rand"
)

type LoginResponse struct {
	Token string `json:"token"`
}

type LoginRequest struct {
	Number   string `json:"number"`
	Password string `json:"password"`
}

type TransferRequest struct {
	ToAccount   string `json:"to_account"`
	FromAccount string `json:"from_account"`
	Amount      int64  `json:"amount"`
}

type CreateAccountRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

type UpdateAccountRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required"`
}

type UpdateAccountPasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

type Account struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Balance   int64  `json:"balance"`
	Number    int64  `json:"number"`
	Password  string `json:"password"`
	Email     string `json:"email"`
}

func ValidatePassword(pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(pw), []byte(pw)) == nil
}

func NewAccount(firstName string, lastName string, email string, password string) (*Account, error) {
	encrypt, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &Account{
		ID:        rand.Intn(10000),
		FirstName: firstName,
		LastName:  lastName,
		Number:    int64(rand.Intn(10000000)),
		Password:  string(encrypt),
		Email:     email,
	}, nil
}
