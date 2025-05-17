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

func NewUserHandler(userController UserController) http.Handler {
	userRouter := mux.NewRouter()

	postRoutes := userRouter.Methods(http.MethodPost).Subrouter()
	postRoutes.Handle("/", userController.Create())

	getRoutes := userRouter.Methods(http.MethodGet).Subrouter()
	getRoutes.Handle("/all", userController.GetAll())
	getRoutes.Handle("/{id:[0-9]+}", userController.GetOne())

	putRoutes := userRouter.Methods(http.MethodPut).Subrouter()
	putRoutes.Handle("/{id:[0-9]+}", userController.Update())

	deleteRoutes := userRouter.Methods(http.MethodDelete).Subrouter()
	deleteRoutes.Handle("/{id:[0-9]+}", userController.Delete())

	return http.Handler(userRouter)
}
