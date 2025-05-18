package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/namph-hanoi/fiddle-golang-restful/util"
)

var testQueries *Queries

var queries Queries

var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("Cannot load config", err)
	}
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannnot connect to the DB. Details: ", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
