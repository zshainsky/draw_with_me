package draw_test

import (
	"fmt"
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	draw "github.com/zshainsky/draw-with-me"
)

type StubSever struct {
	http.Handler
	template *template.Template
}

const htmlTemplatePath = "draw.html"

func TestClient(t *testing.T) {
	t.Run("starting out Client code", func(t *testing.T) {
		got := draw.PrintClient()
		want := "client"

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("create websocket connection", func(t *testing.T) {

		client := draw.NewClient(nil, nil)
		// server := httptest.NewServer(http.HandlerFunc(client.ServeWS))
		stubServer, err := NewStubServer(client)
		if err != nil {
			t.Fatalf("could not start stub server: %v", err)
		}
		server := httptest.NewServer(stubServer)

		url := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"
		fmt.Printf("ws URL: %v\n", url)

		ws, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			t.Fatalf("could not open a ws connection on %s %v", url, err)
		}

		defer server.Close()
		defer ws.Close()

	})

	t.Run("receive broadcast from hub", func(t *testing.T) {

	})

	t.Run("write broadcast to websocket", func(t *testing.T) {

	})
	t.Run("receive messages from websocket", func(t *testing.T) {

	})

	t.Run("send messages to hub's broadcast chanel", func(t *testing.T) {

	})

}

func NewStubServer(c *draw.Client) (*StubSever, error) {
	s := new(StubSever)
	tmpl, err := template.ParseFiles(htmlTemplatePath)

	if err != nil {
		return nil, fmt.Errorf("problem loading template %s", err.Error())
	}
	s.template = tmpl
	fmt.Println("in")

	router := http.NewServeMux()
	router.Handle("/ws", http.HandlerFunc(c.ServeWS))

	s.Handler = router

	return s, nil
}
func writeWSMessage(t *testing.T, conn *websocket.Conn, message string) {
	t.Helper()
	if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		t.Fatalf("could not send message over ws connection %v", err)
	}
}
