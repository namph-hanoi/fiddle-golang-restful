package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/namph-hanoi/fiddle-golang-restful/api"
	db "github.com/namph-hanoi/fiddle-golang-restful/db/sqlc"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
	log.Println("server started at", serverAddress)

}
