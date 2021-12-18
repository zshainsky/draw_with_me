package draw

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"google.golang.org/api/idtoken"
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
				if err != nil { // JWT not valid
					fmt.Println(err)
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte(fmt.Sprintf("Unauthorized: %v", err)))
					panic(err)
				}
				//JWT is valid
				fmt.Printf("\nValid JWT - Claims: %v\n", payload.Claims)
				ctx := context.WithValue(r.Context(), "jwt_payload", payload)
				handler.ServeHTTP(w, r.WithContext(ctx))
			}
		} else {
			// If no token found, redirect to signin page
			fmt.Println(err)
			http.Redirect(w, r, "/signin", http.StatusMovedPermanently)
		}
	})
}
