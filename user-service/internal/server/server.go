package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func Run(router *mux.Router, host string, port int) error {
	var httpServer http.Server = http.Server{
		Addr:         fmt.Sprintf("%s:%d", host, port),
		Handler:      router,
		ErrorLog:     nil,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
		IdleTimeout:  time.Second * 5,
	}
	if err := httpServer.ListenAndServe(); err != nil {
		return err
	}
	return nil
}
