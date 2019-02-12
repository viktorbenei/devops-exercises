package main

import (
	"log"
	"net/http"
	"os"

	"github.com/bitrise-io/api-utils/httpresponse"
	"github.com/pkg/errors"
)

var port = "8182"

func authViaKubernetesSecret(w http.ResponseWriter, r *http.Request) error {
	authSecretToken := os.Getenv("AUTH_SECRET_TOKEN")
	log.Printf("(debug) authSecretToken: %s", authSecretToken)
	if r.Header.Get("Authorization") != "token "+authSecretToken {
		return errors.WithStack(httpresponse.RespondWithUnauthorized(w))
	}
	return errors.WithStack(httpresponse.RespondWithSuccess(w, map[string]string{
		"message": "authorized - OK",
	}))
}

func sayHello(w http.ResponseWriter, r *http.Request) error {
	message := "Hello " + r.URL.Query().Get("name")

	return errors.WithStack(httpresponse.RespondWithSuccess(w, map[string]string{
		"message": message,
		"version": version,
	}))
}

func handleRoot(w http.ResponseWriter, r *http.Request) error {
	return errors.WithStack(httpresponse.RespondWithSuccess(w, map[string]string{
		"message": "Welcome!",
		"version": version,
	}))
}

func main() {
	http.Handle("/auth-via-kubernetes-secret", httpresponse.InternalErrHandlerFuncAdapter(authViaKubernetesSecret))
	http.Handle("/hi", httpresponse.InternalErrHandlerFuncAdapter(sayHello))
	http.Handle("/", httpresponse.InternalErrHandlerFuncAdapter(handleRoot))

	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	log.Printf("Starting (on port: %s) ...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Printf("[!] Exception: %+v", err)
		panic(err)
	}
}
