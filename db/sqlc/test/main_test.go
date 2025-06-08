package db_test

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	db "github.com/namph-hanoi/fiddle-golang-restful/db/sqlc"
	"github.com/namph-hanoi/fiddle-golang-restful/util"
)

var testQueries *db.Queries

var queries db.Queries

var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../../..")
	if err != nil {
		log.Fatal("Cannot load config", err)
	}
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannnot connect to the DB. Details: ", err)
	}

	testQueries = db.New(testDB)

	os.Exit(m.Run())
}
