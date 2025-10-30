package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/datmaithanh/orderfood/utils"
	_ "github.com/lib/pq"
)

var testQueries *Queries

var testDB *sql.DB
var err error

func TestMain(m *testing.M) {
	if err != nil {
		log.Fatal("Can not load config", err)
	}
	testDB, err = sql.Open(utils.DBDriver, utils.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	testQueries = New(testDB)

	code := m.Run()
	os.Exit(code)
}
