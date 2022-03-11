package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nooderg/autotest_templating/config"
	"github.com/nooderg/autotest_templating/core/service"
)

func main() {
	log.Println("Starting server...")
	r := mux.NewRouter()

	log.Println("Connecting to database...")
	config.GetDBClient()

	log.Println("Connected to database!")

	r.HandleFunc("/", service.Generate).Methods("POST")

	log.Println("Server running on localhost:8080!")

	log.Fatal(http.ListenAndServe(":8080", r))
}
