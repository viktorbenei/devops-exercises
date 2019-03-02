package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

type TokenValidator struct {
	hmacSecret []byte
}

func NewTokenValidator(hmacSecret []byte) (*TokenValidator, error) {
	if len(hmacSecret) < 1 {
		return nil, errors.New("Invalid hmac secret (empty)")
	}

	return &TokenValidator{
		hmacSecret: hmacSecret,
	}, nil
}

// Validate ...
func (t *TokenValidator) Validate(r *http.Request) error {
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return errors.New("Only Bearer Authorization accepted")
	}
	jwtToken := strings.TrimPrefix(authHeader, "Bearer ")
	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return t.hmacSecret, nil
	})
	if err != nil {
		return errors.WithStack(err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims["foo"], claims["nbf"])
	} else {
		fmt.Println(err)
	}

	return nil
}
