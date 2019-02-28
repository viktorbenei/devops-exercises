package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

var (
	sessionStore *sessions.FilesystemStore
)

func initAuth0() error {
	sessionStore = sessions.NewFilesystemStore("", []byte("something-very-secret"))
	gob.Register(map[string]interface{}{})
	return nil
}

func auth0CallbackHandler(w http.ResponseWriter, r *http.Request) error {
	domain := os.Getenv("AUTH0_DOMAIN")

	conf := &oauth2.Config{
		ClientID:     os.Getenv("AUTH0_CLIENT_ID"),
		ClientSecret: os.Getenv("AUTH0_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("AUTH0_CALLBACK_URL"),
		Scopes:       []string{"openid", "profile"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://" + domain + "/authorize",
			TokenURL: "https://" + domain + "/oauth/token",
		},
	}

	state := r.URL.Query().Get("state")
	session, err := sessionStore.Get(r, "state")
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
	resp, err := client.Get("https://" + domain + "/userinfo")
	if err != nil {
		return errors.WithStack(err)
	}

	defer resp.Body.Close()

	var profile map[string]interface{}
	if err = json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return errors.WithStack(err)
	}

	session, err = sessionStore.Get(r, "auth-session")
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
	http.Redirect(w, r, "/", http.StatusSeeOther)
	log.Printf("=> token: %#v", token)
	log.Printf("=> session: %#v", session)
	return nil
}

func auth0LoginHandler(w http.ResponseWriter, r *http.Request) error {
	domain := os.Getenv("AUTH0_DOMAIN")
	aud := os.Getenv("AUTH0_AUDIENCE")

	conf := &oauth2.Config{
		ClientID:     os.Getenv("AUTH0_CLIENT_ID"),
		ClientSecret: os.Getenv("AUTH0_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("AUTH0_CALLBACK_URL"),
		Scopes:       []string{"openid", "profile"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://" + domain + "/authorize",
			TokenURL: "https://" + domain + "/oauth/token",
		},
	}

	if aud == "" {
		aud = "https://" + domain + "/userinfo"
	}

	// Generate random state
	b := make([]byte, 32)
	rand.Read(b)
	state := base64.StdEncoding.EncodeToString(b)

	session, err := sessionStore.Get(r, "state")
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
