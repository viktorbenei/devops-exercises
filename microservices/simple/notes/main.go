package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/bitrise-io/api-utils/httpresponse"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

const appName = "private.notes"

var port = "8182"
var datastore *Datastore

func handleGetNotes(w http.ResponseWriter, r *http.Request) error {
	userID := r.Header.Get("UserID")

	userNotes, err := datastore.GetNotes(UserID(userID))
	if err != nil {
		return errors.WithStack(err)
	}

	return errors.WithStack(httpresponse.RespondWithSuccess(w, map[string]map[NoteID]Note{
		"notes": userNotes,
	}))
}

func handleCreateNote(w http.ResponseWriter, r *http.Request) error {
	userID := r.Header.Get("UserID")
	noteID := r.Header.Get("NoteID")

	note := Note{}
	defer httpresponse.RequestBodyCloseWithErrorLog(r)
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		return httpresponse.RespondWithBadRequestError(w, "Invalid request body, JSON decode failed")
	}

	if err := datastore.SetNote(UserID(userID), NoteID(noteID), note); err != nil {
		return errors.WithStack(err)
	}

	return errors.WithStack(httpresponse.RespondWithSuccess(w, map[string]Note{
		"note": note,
	}))
}

func handleRoot(w http.ResponseWriter, r *http.Request) error {
	return errors.WithStack(httpresponse.RespondWithSuccess(w, map[string]string{
		"message": "Welcome!",
		"app":     appName,
		"version": version,
	}))
}

func handlerErrorHandlingAdapter(h httpresponse.HanderFuncWithInternalError) http.Handler {
	// TODO: instead of redefining httpresponse.InternalErrHandlerFuncAdapter extend it
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		intServErr := h(w, r)
		if intServErr != nil {
			if inputErr, ok := errors.Cause(intServErr).(*InputError); ok {
				httpresponse.RespondWithBadRequestErrorNoErr(w, inputErr.Error())
			} else if notFoundErr, ok := errors.Cause(intServErr).(*NotFoundError); ok {
				httpresponse.RespondWithErrorNoErr(w, notFoundErr.Error(), http.StatusNotFound)
			} else {
				httpresponse.RespondWithInternalServerError(w, errors.WithStack(intServErr))
			}
		}
	})
}

func mainE() error {
	// Init
	datastore = NewDatastore()

	// Setup routing
	r := mux.NewRouter().StrictSlash(true)
	middlewareProvider := NewMiddlewareProvider(appName, version)
	//
	r.Handle("/notes", middlewareProvider.CommonMiddleware().Then(
		handlerErrorHandlingAdapter(handleGetNotes))).Methods("GET")
	r.Handle("/notes", middlewareProvider.CommonMiddleware().Then(
		handlerErrorHandlingAdapter(handleCreateNote))).Methods("POST")
	r.Handle("/", middlewareProvider.CommonMiddleware().Then(
		handlerErrorHandlingAdapter(handleRoot))).Methods("GET")
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
