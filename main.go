package main

import (
	"github.com/mervanerdem/PropertyFinderSaleAPI/Server"
	"github.com/mervanerdem/PropertyFinderSaleAPI/SqlConnection"
	"log"
	"net/http"
)

func main() {

	storage, _, err := SqlConnection.NewMStorage("root:Mervan.1907@tcp(127.0.0.1:3306)/pfsale")
	if err != nil {
		panic("Configration is wrong")
	}
	log.Fatal(http.ListenAndServe("localhost:8080", Server.NewServer(storage)))
}
