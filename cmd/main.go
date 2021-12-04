package main

import (
	"fmt"
	"log"
	"net/http"

	draw "github.com/zshainsky/draw-with-me"
)

const htmlFileName = "../draw.html"

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, htmlFileName)
}

func main() {
	fmt.Println(draw.PrintHub())
	hub := draw.NewHub()
	go hub.Run()

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		draw.ServeWS(hub, w, r)
	})

	log.Fatal(http.ListenAndServe(":3000", nil))

}
