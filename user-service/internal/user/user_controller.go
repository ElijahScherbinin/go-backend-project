package user

import (
	"net/http"

	"github.com/gorilla/mux"
)

type UserController interface {
	Create() http.Handler
	GetAll() http.Handler
	GetOne() http.Handler
	Update() http.Handler
	Delete() http.Handler
}

func SetupUserRoutes(router *mux.Router, userController UserController) *mux.Router {

	postRoutes := router.Methods(http.MethodPost).Subrouter()
	postRoutes.Handle("/users", userController.Create())

	getRoutes := router.Methods(http.MethodGet).Subrouter()
	getRoutes.Handle("/users/all", userController.GetAll())
	getRoutes.Handle("/users/{id:[0-9]+}", userController.GetOne())

	putRoutes := router.Methods(http.MethodPut).Subrouter()
	putRoutes.Handle("/users/{id:[0-9]+}", userController.Update())

	deleteRoutes := router.Methods(http.MethodDelete).Subrouter()
	deleteRoutes.Handle("/users/{id:[0-9]+}", userController.Delete())

	return router
}
