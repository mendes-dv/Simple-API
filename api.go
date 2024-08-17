package main

import (
	"awesomeProject/models"
	"encoding/json"
	"fmt"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"strconv"
)

type APIServer struct {
	listenAddr string
	store      Storage
}

func NewAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/login", makeHttpHandleFunc(s.handleLogin))

	router.HandleFunc("/account", makeHttpHandleFunc(s.handleAccount))

	router.HandleFunc("/account/{id}", withJWTAuth(makeHttpHandleFunc(s.handleAccount), s.store))

	router.HandleFunc("/transfer", makeHttpHandleFunc(s.handleTransfer))

	log.Println("Listening on port:", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "POST" {
		var req models.LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return err
		}

		number, err := strconv.Atoi(req.Number)

		acc, err := s.store.GetAccountByNumber(number)
		if err != nil {
			return err
		}

		if !models.ValidatePassword(req.Number) {
			return fmt.Errorf("Invalid password")
		}

		token, err := createJWT(acc)
		if err != nil {
			return err
		}

		resp := models.LoginResponse{
			Token: token,
		}

		return WriteJson(w, http.StatusOK, resp)
	}

	return WriteJson(w, http.StatusMethodNotAllowed, "method not allowed")
}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]

	if r.URL.Path == "/account" && id == "" {
		if r.Method == http.MethodGet {
			return s.handleGetAccounts(w, r)
		}
		if r.Method == http.MethodPost {
			return s.handleCreateAccount(w, r)
		}
	}

	if id != "" {
		if r.Method == http.MethodGet {
			return s.handleGetAccountById(w, r)
		}
		if r.Method == http.MethodDelete {
			return s.handleDeleteAccount(w, r)
		}
	}

	return WriteJson(w, http.StatusMethodNotAllowed, fmt.Errorf("unsupported method: %s", r.Method))
}

func (s *APIServer) handleGetAccounts(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.store.GetAccounts()
	if err != nil {
		return err
	}

	return WriteJson(w, http.StatusOK, accounts)
}

func (s *APIServer) handleGetAccountById(w http.ResponseWriter, r *http.Request) error {
	id, err := getId(r)
	if err != nil {
		return err
	}

	account, err := s.store.GetAccountById(id)
	if err != nil {
		return err
	}

	return WriteJson(w, http.StatusOK, account)
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	createAccountRequest := new(models.CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(createAccountRequest); err != nil {
		return err
	}

	account, _ := models.NewAccount(createAccountRequest.FirstName, createAccountRequest.LastName, createAccountRequest.Email, createAccountRequest.Password)
	err := s.store.CreateAccount(account)
	if err != nil {
		return err
	}

	return WriteJson(w, http.StatusOK, account)
}

func createJWT(account *models.Account) (string, error) {
	claims := jwt.MapClaims{
		"expiresAt":     15000,
		"accountNumber": account.Number,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := os.Getenv("SECRET_KEY")
	return token.SignedString([]byte(secretKey)) //
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := getId(r)
	if err != nil {
		return err
	}

	if err := s.store.DeleteAccount(id); err != nil {
		return err
	}

	return WriteJson(w, http.StatusOK, "Account was deleted")
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	transferRequest := new(models.TransferRequest)
	if err := json.NewDecoder(r.Body).Decode(transferRequest); err != nil {
		return err
	}
	defer r.Body.Close()

	return WriteJson(w, http.StatusOK, transferRequest)
}

func WriteJson(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func withJWTAuth(handlerFunc http.HandlerFunc, store Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("calling JWT middleware")
		tokenString := r.Header.Get("Authorization")
		token, err := validateJWT(tokenString)
		if err != nil {
			WriteJson(w, http.StatusUnauthorized, APIError{Error: "invalid token"})
			return
		}

		if !token.Valid {
			WriteJson(w, http.StatusUnauthorized, APIError{Error: "invalid token"})
		}

		userId, err := getId(r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		account, err := store.GetAccountById(userId)
		if err != nil {
			WriteJson(w, http.StatusUnauthorized, APIError{Error: "invalid token"})
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		if account.Number != int64(claims["accountNumber"].(float64)) {
			WriteJson(w, http.StatusUnauthorized, APIError{Error: "invalid token"})
			return
		}

		handlerFunc(w, r)
	}
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
}

type APIFunc func(w http.ResponseWriter, r *http.Request) error

type APIError struct {
	Error string `json:"error"`
}

func makeHttpHandleFunc(fn APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			WriteJson(w, http.StatusBadRequest, APIError{Error: err.Error()})
		}
	}
}

func getId(r *http.Request) (int, error) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("invalid id: %s", idStr)
	}

	return id, nil
}
