package routes

import (
	"github.com/gorilla/mux"
	"github.com/yash04g/Go-mux-mongo/controllers"
)

func UserRoute(router *mux.Router) {
	router.HandleFunc("/user", controllers.CreateUser()).Methods("POST")
	router.HandleFunc("/user/{userId}", controllers.GetAUser()).Methods("GET")
	router.HandleFunc("/user/{userId}", controllers.EditAUser()).Methods("PUT")
	router.HandleFunc("/user/{userId}", controllers.DeleteAUser()).Methods("DELETE")
	router.HandleFunc("/users", controllers.GetAllUsers()).Methods("GET")
}
