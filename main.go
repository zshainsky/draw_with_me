package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/zshainsky/draw-with-me/draw"
)

var router *mux.Router

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
