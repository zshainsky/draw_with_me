package main

import (
	"log"
	"net/http"

	"github.com/zshainsky/draw-with-me"
)

var rooms []*draw.Room

const htmlFileName = "../home.html"

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
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/create-room", draw.ServeRoom)

	log.Fatal(http.ListenAndServe(":3000", nil))

}
