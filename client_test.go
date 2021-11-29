package draw_test

import (
	"testing"

	draw "github.com/zshainsky/draw-with-me"
)

func TestClient(t *testing.T) {
	t.Run("starting out Client code", func(t *testing.T) {
		got := draw.PrintClient()
		want := "client"

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("create websocket connection", func(t *testing.T) {
		// client := draw.NewClient()
		// s := httptest.NewServer(http.HandlerFunc())
	})

}
