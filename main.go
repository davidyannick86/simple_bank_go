package main

import (
	"database/sql"
	"log"

	"github.com/davidyannick86/simple-bank/api"
	db "github.com/davidyannick86/simple-bank/db/sqlc"
	"github.com/davidyannick86/simple-bank/utils"
	_ "github.com/lib/pq"
)

func main() {

	config, err := utils.LoadConfig(".")

	conn, err := sql.Open(config.DbDriver, config.DbSource)
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
