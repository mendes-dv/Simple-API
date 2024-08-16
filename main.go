package main

import (
	"log"
)

func main() {
	postgresStore, err := NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}

	if err := postgresStore.Init(); err != nil {
		log.Fatal(err)
	}

	server := NewAPIServer(":3000", postgresStore)
	server.Run()
}
