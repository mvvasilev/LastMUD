package auth

import (
	"code.haedhutner.dev/mvv/LastMUD/services/auth/internal/handlers"
	"code.haedhutner.dev/mvv/LastMUD/shared/httputils"
	"code.haedhutner.dev/mvv/LastMUD/shared/log"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

func LaunchAuthServiceServer() {
	logger := log.NewLogger("auth_service", log.DBG, log.DBG, "./log/auth_service.log")
	r := mux.NewRouter()

	r.Use(httputils.NewLoggingMiddleware(logger))
	r.Use(httputils.NewErrorHandlingMiddleware())

	api := r.PathPrefix("/api/v1/account").Subrouter()

	api.Methods(http.MethodGet).Path("/{accountId}").Handler(handlers.NewGetAccountHandler())
	api.Methods(http.MethodPost).Path("/").Handler(handlers.NewCreateAccountHandler())
	api.Methods(http.MethodPut).Path("/{accountId}").Handler(handlers.NewUpdateAccountHandler())
	api.Methods(http.MethodDelete).Path("/{accountId}").Handler(handlers.NewDeleteAccountHandler())

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8001",
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
