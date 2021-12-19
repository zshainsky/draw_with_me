package draw

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"google.golang.org/api/idtoken"
)

type CTXKey string

const (
	jwtCTXKey CTXKey = "jwt"
)

// Created incase add additional auth types
type AuthType string

const (
	google AuthType = "google"
)

const googleClientId = "406504108908-4djtjr6q3lil4rgrnbjproqi7ruc59vs.apps.googleusercontent.com"

func AuthMiddleware(handler http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("\n Auth Handler: r.Header: %+v \n", r.Header)
		if jwtCookie, err := r.Cookie("jwt-token"); jwtCookie != nil {
			token := strings.Split(jwtCookie.Value, " ")
			//index [0] should be the word "Bearer" and index [1] should be the token value
			if len(token) != 2 {
				fmt.Println("Malformed token")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Malformed Token"))
			} else {
				// Validate JWT
				jwt := token[1]
				payload, err := idtoken.Validate(context.Background(), jwt, googleClientId)
				if err != nil { // JWT not valid - redirect to signin screen
					fmt.Println(err)
					http.Redirect(w, r, "/signin", http.StatusMovedPermanently)
					return
				}
				//JWT is valid
				fmt.Printf("\nValid JWT - Claims: %v\n", payload.Claims)
				// Check if this is an authorization request with POST method from Auth provider. Create user if so
				if r.URL.Path == "/authorize" && r.Method == "POST" {
					//Does the user exist?

					//If not, create
				}
				ctx := context.WithValue(r.Context(), jwtCTXKey, payload)
				handler.ServeHTTP(w, r.WithContext(ctx))
			}
		} else {
			// If no token found, redirect to signin page
			fmt.Println(err)
			http.Redirect(w, r, "/signin", http.StatusMovedPermanently)
		}
	})
}
