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
	DeleteAccount(id int) error
	UpdateAccount(account *models.Account) error
	GetAccounts() ([]*models.Account, error)
	GetAccountById(id int) (*models.Account, error)
	GetAccountByNumber(number int) (*models.Account, error)
}

type PostgresSQLStorage struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresSQLStorage, error) {
	connStr := "user=postgres dbname=postgres password=gobank sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return &PostgresSQLStorage{db: db}, nil
}

func (s *PostgresSQLStorage) Init() error {
	return s.createAccountTable()
}

func (s *PostgresSQLStorage) createAccountTable() error {
	log.Println("Creating account table")
	query := `CREATE TABLE IF NOT EXISTS accounts (
		id SERIAL PRIMARY KEY,
		first_name TEXT,
		last_name TEXT,
		balance BIGINT,
		number BIGINT,
		email TEXT,
		password TEXT,
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

func (s *PostgresSQLStorage) CreateAccount(account *models.Account) error {
	query := `INSERT INTO accounts 
	(first_name, last_name, balance, number, email, password)
	VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := s.db.Exec(query, account.FirstName, account.LastName, account.Balance, account.Number, account.Email, account.Password)
	if err != nil {
		log.Println("Error creating account:", err)
		return err
	}

	return nil
}

func (s *PostgresSQLStorage) GetAccountByNumber(number int) (*models.Account, error) {
	query := `
	SELECT id, first_name, last_name, balance, number, email 
	FROM accounts
	WHERE number = $1`
	rows, err := s.db.Query(query, number)
	if err != nil {
		log.Println("Error retrieving account by number:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		return scanIntoAccount(rows)
	}

	return nil, fmt.Errorf("account with number %d not found", number)
}

func (s *PostgresSQLStorage) DeleteAccount(id int) error {
	query := "DELETE FROM accounts WHERE id=$1"
	_, err := s.db.Exec(query, id)
	if err != nil {
		log.Println("Error deleting account:", err)
		return err
	}
	return nil
}

func (s *PostgresSQLStorage) UpdateAccount(account *models.Account) error {
	query := `UPDATE accounts SET 
	first_name=$1, last_name=$2, balance=$3, number=$4, email=$5, password=$6, updated_at=CURRENT_TIMESTAMP 
	WHERE id=$7`

	_, err := s.db.Exec(query, account.FirstName, account.LastName, account.Balance, account.Number, account.Email, account.Password, account.ID)
	if err != nil {
		log.Println("Error updating account:", err)
		return err
	}

	return nil
}

func (s *PostgresSQLStorage) GetAccountById(id int) (*models.Account, error) {
	query := `
	SELECT id, first_name, last_name, balance, number, email
	FROM accounts
	WHERE id = $1`
	rows, err := s.db.Query(query, id)
	if err != nil {
		log.Println("Error retrieving account by ID:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		return scanIntoAccount(rows)
	}

	return nil, fmt.Errorf("account %d not found", id)
}

func (s *PostgresSQLStorage) GetAccounts() ([]*models.Account, error) {
	rows, err := s.db.Query("SELECT id, first_name, last_name, balance, number, email FROM accounts")
	if err != nil {
		log.Println("Error retrieving accounts:", err)
		return nil, err
	}
	defer rows.Close()

	var accounts []*models.Account
	for rows.Next() {
		account, err := scanIntoAccount(rows)
		if err != nil {
			log.Println("Error scanning account:", err)
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func scanIntoAccount(rows *sql.Rows) (*models.Account, error) {
	account := &models.Account{}
	err := rows.Scan(&account.ID, &account.FirstName, &account.LastName, &account.Balance, &account.Number, &account.Email, &account.Password)
	return account, err
}
