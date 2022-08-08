package main

import (
	"database/sql"
	"github.com/mervanerdem/PropertyFinderSaleAPI/server"
	"github.com/mervanerdem/PropertyFinderSaleAPI/services"
	"github.com/mervanerdem/PropertyFinderSaleAPI/sqlConnection"
	"log"
	"net/http"
)

func main() {
	dsn := services.GetDsn()
	storage, db, err := sqlConnection.NewMStorage(dsn)
	if err != nil {
		log.Fatal("Configuration is wrong")
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)
	hostAddress := services.GetHost()
	log.Fatal(http.ListenAndServe(hostAddress, server.NewServer(storage)))

}
