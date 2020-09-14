package router

import (
	"go-postgres-crud/middleware"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	routes := mux.NewRouter().StrictSlash(true)

	routes.HandleFunc("/api/user", middleware.GetAllUser).Methods("GET", "OPTIONS")
	routes.HandleFunc("/api/newuser", middleware.CreateUser).Methods("POST", "OPTIONS")
	routes.HandleFunc("/api/user/{id}", middleware.GetUser).Methods("GET", "OPTIONS")

	return routes

}
