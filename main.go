package main

import (
	"awesomeProject/models"
	"flag"
	"fmt"
	"log"
)

func seedAccount(store Storage, fname, lname, email, pw string) *models.Account {
	acc, err := models.NewAccount(fname, lname, email, pw)
	if err != nil {
		log.Fatal(err)
	}

	if err := store.CreateAccount(acc); err != nil {
		log.Fatal(err)
	}

	return acc
}

func seedAccounts(s Storage) {
	seedAccount(s, "12313", "Monkey", "teste@mail.com", "112233")
	seedAccount(s, "monkey", "Jonh", "monkey@mail.com", "223344")
}

type PostgresStore struct {
	// Add required store fields, e.g., DB connection
}

func main() {
	seed := flag.Bool("s", false, "seed")
	flag.Parse()

	postgresStore, err := NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}

	if err := postgresStore.Init(); err != nil {
		log.Fatal(err)
	}

	if *seed {
		fmt.Println("seeding ...")
		seedAccounts(postgresStore)
	}

	server := NewAPIServer(":3000", postgresStore)
	server.Run()
}
