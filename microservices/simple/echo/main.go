package main

import (
	"log"
	"net/http"
	"os"

	"github.com/bitrise-io/api-utils/httpresponse"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

const appName = "echo.api"

var port = "8182"

func sayHello(w http.ResponseWriter, r *http.Request) error {
	message := "Hello " + r.URL.Query().Get("name")

	return errors.WithStack(httpresponse.RespondWithSuccess(w, map[string]string{
		"message": message,
	}))
}

func handleRoot(w http.ResponseWriter, r *http.Request) error {
	return errors.WithStack(httpresponse.RespondWithSuccess(w, map[string]string{
		"message": "Welcome!",
		"app":     appName,
		"version": version,
	}))
}

func mainE() error {
	// Setup routing
	r := mux.NewRouter().StrictSlash(true)
	middlewareProvider := NewMiddlewareProvider(appName, version)
	//
	r.Handle("/hi", middlewareProvider.CommonMiddleware().Then(
		httpresponse.InternalErrHandlerFuncAdapter(sayHello))).Methods("GET")
	r.Handle("/", middlewareProvider.CommonMiddleware().Then(
		httpresponse.InternalErrHandlerFuncAdapter(handleRoot))).Methods("GET")
	//
	http.Handle("/", r)

	// Start the server
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}
	log.Printf("Starting (on port: %s) ...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Printf("[!] Exception: %+v", err)
		panic(err)
	}
	return nil
}

func main() {
	if err := mainE(); err != nil {
		log.Printf("[!] Exception: %+v", err)
	}
}