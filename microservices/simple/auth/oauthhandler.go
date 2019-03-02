package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

// OAuthHandler ...
type OAuthHandler struct {
	domain       string
	clientID     string
	clientSecret string
	callbackURL  string
	//
	audience string

	//
	sessionStore sessions.Store
}

func requiredEnv(envKey string) (string, error) {
	// TODO: replace this with https://github.com/bitrise-io/go-utils/pull/85 once that's merged
	val := os.Getenv(envKey)
	if len(val) < 1 {
		return "", errors.Errorf("Required environment variable (%s) not provided", envKey)
	}
	return val, nil
}

// NewOAuthHandler ...
func NewOAuthHandler(sessionStore sessions.Store) (*OAuthHandler, error) {
	domain, err := requiredEnv("OAUTH_DOMAIN")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	clientID, err := requiredEnv("OAUTH_CLIENT_ID")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	clientSecret, err := requiredEnv("OAUTH_CLIENT_SECRET")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	callbackURL, err := requiredEnv("OAUTH_CALLBACK_URL")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	audience := os.Getenv("OAUTH_AUDIENCE")
	if len(audience) < 1 {
		audience = "https://" + domain + "/userinfo"
		log.Printf(" (!) OAUTH_AUDIENCE was not provided, using the domain instead: %s", audience)
	}

	return &OAuthHandler{
		domain:       domain,
		clientID:     clientID,
		clientSecret: clientSecret,
		callbackURL:  callbackURL,
		//
		audience: audience,
		//
		sessionStore: sessionStore,
	}, nil
}

// OAuthCallbackHandler ...
func (ah *OAuthHandler) OAuthCallbackHandler(w http.ResponseWriter, r *http.Request) error {
	conf := &oauth2.Config{
		ClientID:     ah.clientID,
		ClientSecret: ah.clientSecret,
		RedirectURL:  ah.callbackURL,
		Scopes:       []string{"openid", "profile"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://" + ah.domain + "/authorize",
			TokenURL: "https://" + ah.domain + "/oauth/token",
		},
	}

	state := r.URL.Query().Get("state")
	session, err := ah.sessionStore.Get(r, "state")
	if err != nil {
		return errors.WithStack(err)
	}

	if state != session.Values["state"] {
		http.Error(w, "Invalid state parameter", http.StatusInternalServerError)
		return nil
	}

	code := r.URL.Query().Get("code")

	token, err := conf.Exchange(context.TODO(), code)
	if err != nil {
		return errors.WithStack(err)
	}

	// Getting now the userInfo
	client := conf.Client(context.TODO(), token)
	resp, err := client.Get("https://" + ah.domain + "/userinfo")
	if err != nil {
		return errors.WithStack(err)
	}

	defer resp.Body.Close()

	var profile map[string]interface{}
	if err = json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return errors.WithStack(err)
	}

	session, err = ah.sessionStore.Get(r, "auth-session")
	if err != nil {
		return errors.WithStack(err)
	}

	session.Values["id_token"] = token.Extra("id_token")
	session.Values["access_token"] = token.AccessToken
	session.Values["profile"] = profile
	err = session.Save(r, w)
	if err != nil {
		return errors.WithStack(err)
	}

	// Redirect to logged in page
	http.Redirect(w, r, "/token/user", http.StatusSeeOther)
	return nil
}

// OAuthLoginHandler ...
func (ah *OAuthHandler) OAuthLoginHandler(w http.ResponseWriter, r *http.Request) error {
	aud := ah.audience

	conf := &oauth2.Config{
		ClientID:     ah.clientID,
		ClientSecret: ah.clientSecret,
		RedirectURL:  ah.callbackURL,
		Scopes:       []string{"openid", "profile"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://" + ah.domain + "/authorize",
			TokenURL: "https://" + ah.domain + "/oauth/token",
		},
	}

	// Generate random state
	b := make([]byte, 32)
	rand.Read(b)
	state := base64.StdEncoding.EncodeToString(b)

	session, err := ah.sessionStore.Get(r, "state")
	if err != nil {
		return errors.WithStack(err)
	}
	session.Values["state"] = state
	err = session.Save(r, w)
	if err != nil {
		return errors.WithStack(err)
	}

	audience := oauth2.SetAuthURLParam("audience", aud)
	url := conf.AuthCodeURL(state, audience)

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	return nil
}
