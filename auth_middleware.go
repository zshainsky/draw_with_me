package draw

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

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
		// fmt.Printf("\n Auth Handler: r.Header: %+v \n", r.Header)
		fmt.Printf("\n Auth Handler running from route: %v\n", r.URL)
		if jwtCookie, err := r.Cookie("jwt-token"); jwtCookie != nil {
			fmt.Printf("\njwt-token found")
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
				// payload, err := idtoken.Validate(context.Background(),)
				if err != nil { // JWT not valid - redirect to signin screen
					fmt.Println(err)
					// TODO: Untested...potentially introduce this if the validator with credentials doesn't work to refresh token
					// this might throw nil pointer exception...because payload may be nil
					if time.Now().Unix() > payload.Expires {
						// Refresh token
						fmt.Printf("idtoken: token expired...refreshing token")
						// return
					}
					http.Redirect(w, r, "/signin", http.StatusMovedPermanently)
					return
				}
				//JWT is valid
				fmt.Printf("\nValid JWT - Claims: %v\n", payload.Claims)

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

func now() {
	panic("unimplemented")
}
