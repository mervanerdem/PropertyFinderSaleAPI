package main

import (
	"github.com/mervanerdem/PropertyFinderSaleAPI/server"
	"github.com/mervanerdem/PropertyFinderSaleAPI/sqlConnection"
	"log"
	"net/http"
)

func main() {

	storage, db, err := sqlConnection.NewMStorage("propertyfinder:password@tcp(localhost:3306)/pfsale")
	if err != nil {
		log.Fatal("Configuration is wrong")
	}
	defer db.Close()
	log.Fatal(http.ListenAndServe("localhost:8080", server.NewServer(storage)))

}
