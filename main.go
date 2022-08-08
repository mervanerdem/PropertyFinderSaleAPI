package main

import (
	"github.com/mervanerdem/PropertyFinderSaleAPI/Server"
	"github.com/mervanerdem/PropertyFinderSaleAPI/SqlConnection"
	"log"
	"net/http"
)

func main() {

	storage, db, err := SqlConnection.NewMStorage("propertyfinder:password@tcp(localhost:3306)/pfsale")
	if err != nil {
		log.Fatal("Configuration is wrong")
	}
	defer db.Close()
	log.Fatal(http.ListenAndServe("localhost:8080", Server.NewServer(storage)))
}
