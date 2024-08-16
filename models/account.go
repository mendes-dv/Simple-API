package models

import (
	"math/rand"
)

type TransferRequest struct {
	ToAccount   string `json:"to_account"`
	FromAccount string `json:"from_account"`
	Amount      int64  `json:"amount"`
}

type CreateAccountRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required"`
}

type Account struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Balance   int64  `json:"balance"`
	Number    int64  `json:"number"`
	Email     string `json:"email"`
}

func NewAccount(firstName string, lastName string, email string) *Account {
	return &Account{
		ID:        rand.Intn(10000),
		FirstName: firstName,
		LastName:  lastName,
		Number:    int64(rand.Intn(10000000)),
		Email:     email,
	}
}
