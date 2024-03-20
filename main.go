package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"simple-bank/api"
	db "simple-bank/db/sqlc"
	until "simple-bank/util"
)

func main() {
	config, err := until.LoadConfig(".")
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	store := db.NewStore(conn)
	server := api.NewServer(store)
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
