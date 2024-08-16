package main

import (
	"awesomeProject/models"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

type Storage interface {
	CreateAccount(account *models.Account) error
	DeleteAccount(int) error
	UpdateAccount(account *models.Account) error
	GetAccountById(int) (*models.Account, error)
	GetAccounts() ([]*models.Account, error)
}

type PostgreSQLStorage struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgreSQLStorage, error) {
	connStr := "user=postgres dbname=postgres password=gobank sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return &PostgreSQLStorage{db: db}, nil
}

func (s *PostgreSQLStorage) Init() error {
	return s.createAccountTable()
}
func (s *PostgreSQLStorage) createAccountTable() error {
	log.Println("Creating account table")
	query := `CREATE TABLE IF NOT EXISTS accounts (
		id SERIAL PRIMARY KEY,
		first_name TEXT,
		last_name TEXT,
		balance BIGINT,
		number BIGINT,
		email TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`
	_, err := s.db.Exec(query)
	if err != nil {
		log.Println("Error creating account table:", err)
		return err
	}
	return nil
}

func (s *PostgreSQLStorage) CreateAccount(account *models.Account) error {
	query := ` insert into accounts 
	( first_name,last_name, balance, number, email)
	VALUES ($1, $2, $3, $4, $5)`

	resp, err := s.db.Exec(query, account.FirstName, account.LastName, account.Balance, account.Number, account.Email)
	if err != nil {
		return err
	}

	fmt.Println(resp)

	return nil
}

func (s *PostgreSQLStorage) DeleteAccount(int) error {
	return nil
}

func (s *PostgreSQLStorage) UpdateAccount(account *models.Account) error {
	return nil
}

func (s *PostgreSQLStorage) GetAccountById(int) (*models.Account, error) {
	return nil, nil
}

func (s *PostgreSQLStorage) GetAccounts() ([]*models.Account, error) {

	rows, err := s.db.Query("SELECT id, first_name, last_name, balance, number, email FROM accounts")

	if err != nil {
		return nil, err
	}

	accounts := make([]*models.Account, 0)

	for rows.Next() {
		account := &models.Account{}
		err = rows.Scan(&account.Id, &account.FirstName, &account.LastName, &account.Balance, &account.Number, &account.Email)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}
