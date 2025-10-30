package main

import (
	"database/sql"
	"log"

	"github.com/datmaithanh/orderfood/api"
	db "github.com/datmaithanh/orderfood/db/sqlc"
	"github.com/datmaithanh/orderfood/utils"
	_ "github.com/lib/pq"
)

func main() {
	conn, err := sql.Open(utils.DBDriver, utils.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(store)
	if err != nil {
		log.Fatal("Cannot run server: %w", err)
	}
	err = server.Start(utils.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
