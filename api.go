package main

import (
	"awesomeProject/models"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
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

	router.HandleFunc("/account", makeHttpHandleFunc(s.handleAccount))

	router.HandleFunc("/account/{id}", makeHttpHandleFunc(s.handleAccount))

	router.HandleFunc("/transfer", makeHttpHandleFunc(s.handleTransfer))

	log.Println("Listening on port:", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)
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

	account := models.NewAccount(createAccountRequest.FirstName, createAccountRequest.LastName, createAccountRequest.Email)
	err := s.store.CreateAccount(account)
	if err != nil {
		return err
	}

	return WriteJson(w, http.StatusOK, account)
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
