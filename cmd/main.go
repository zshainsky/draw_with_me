package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/zshainsky/draw-with-me"
)

var router *mux.Router

type AuthRequestBody struct {
	Credential   string `json:credential`
	G_csrf_token string `json:g_csrf_token`
}

type APIResponse struct {
	Code int //should be http.<response code>
}

func main() {

	router = mux.NewRouter()
	draw.NewServer(router)
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000" // Default port if not specified
	}

	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		log.Fatal(err)
	}

}
