package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yash04g/Go-mux-mongo/configs"
	"github.com/yash04g/Go-mux-mongo/routes"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"data": "Hello World! Welcome to Mongo and Mux"})
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", HomePage).Methods("GET")
	// Run database
	configs.ConnectDB()

	routes.UserRoute(router)

	log.Fatal(http.ListenAndServe(":8080", router))
}
