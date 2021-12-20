package main

import (
	"log"
	"net/http"

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

	err := http.ListenAndServe(":3000", router)
	if err != nil {
		log.Fatal(err)
	}

}
