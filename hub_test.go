package draw_test

import (
	"reflect"
	"testing"

	draw "github.com/zshainsky/draw-with-me"
)

func TestHUB(t *testing.T) {
	t.Run("starting out HUB code", func(t *testing.T) {
		got := draw.PrintHub()
		want := "Hub"

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("successfully add client to registered list", func(t *testing.T) {
		hub := draw.NewHub()
		client := draw.NewClient(nil)

		hub.RegisterClient(client)
		got, err := hub.GetClient(client)
		want := client
		assertRegisteredClientId(t, got, want, err)
	})

	t.Run("remove client from registered list", func(t *testing.T) {
		hub := draw.NewHub()
		client := draw.NewClient(nil)

		// First register client and assert it worked correctly
		hub.RegisterClient(client)
		got, err := hub.GetClient(client)
		want := client
		assertRegisteredClientId(t, got, want, err)

		// Test unregistering the client
		hub.UnregisterClient(client)
		got, _ = hub.GetClient(client)
		assertUnregisteredClient(t, got)

	})

	t.Run("register client", func(t *testing.T) {

	})
}

func assertRegisteredClientId(t testing.TB, got, want *draw.Client, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("asserting: Hub could not find client, err: %v", err)
	}

	if !reflect.DeepEqual(got.GetId(), want.GetId()) {
		t.Errorf("asserting: Hub did not add client successfully. Got: %v, Want: %v", got, want)
	}

}

func assertUnregisteredClient(t testing.TB, got *draw.Client) {
	t.Helper()
	if got != nil {
		t.Errorf("client should have been UNREGISTERED, %v", got.GetId())
	}
}
