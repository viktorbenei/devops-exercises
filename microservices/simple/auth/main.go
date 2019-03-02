package main

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bitrise-io/api-utils/httpresponse"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/pkg/errors"
)

const appName = "public.auth"

var port = "8182"

var globalSessionStore sessions.Store
var globalJWTHMACSecret []byte

func handleGenUserJWT(w http.ResponseWriter, r *http.Request) error {
	session, err := globalSessionStore.Get(r, "auth-session")
	if err != nil {
		return errors.WithStack(err)
	}

	idToken := session.Values["id_token"]
	profile := session.Values["profile"]
	log.Printf("=> id_token: %#v", idToken)
	log.Printf("=> access_token: %#v", session.Values["access_token"])
	log.Printf("=> profile: %#v", profile)
	log.Printf("=> session: %#v", session)

	// --------------------------------------------------------------------------------
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	profileCasted, ok := (profile).(map[string]interface{})
	if !ok {
		return errors.New("Failed to cast 'profile' to required type")
	}
	sub, ok := (profileCasted["sub"]).(string)
	if !ok {
		return errors.New("Failed to cast 'sub' to required type")
	}
	claims := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(10 * time.Minute).Unix(),
		IssuedAt:  time.Now().Unix(),
		Issuer:    appName,
		Subject:   sub,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(globalJWTHMACSecret)

	return errors.WithStack(httpresponse.RespondWithSuccess(w, map[string]interface{}{
		"jwt": tokenString,
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
	// Init
	//  Session
	{
		// TODO: either remove the session store or handle it properly.
		sessionKey, err := requiredEnv("SESSION_KEY")
		if err != nil {
			return errors.WithStack(err)
		}
		globalSessionStore = sessions.NewCookieStore([]byte(sessionKey))
		gob.Register(map[string]interface{}{})
	}
	//  OAuth handler
	authHandler, err := NewOAuthHandler(globalSessionStore)
	if err != nil {
		return errors.WithStack(err)
	}
	//
	{
		jwtHmacSecret, err := requiredEnv("JWT_HMAC_SECRET")
		if err != nil {
			return errors.WithStack(err)
		}
		globalJWTHMACSecret = []byte(jwtHmacSecret)
	}

	// Setup routing
	r := mux.NewRouter().StrictSlash(true)
	middlewareProvider := NewMiddlewareProvider(appName, version)
	// auth
	r.Handle("/callback", httpresponse.InternalErrHandlerFuncAdapter(authHandler.OAuthCallbackHandler))
	r.Handle("/login/user", httpresponse.InternalErrHandlerFuncAdapter(authHandler.OAuthLoginHandler))
	r.Handle("/token/user", httpresponse.InternalErrHandlerFuncAdapter(handleGenUserJWT))
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
