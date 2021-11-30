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
		client := draw.NewClient(nil, hub)

		hub.RegisterClient(client)
		got, err := hub.GetClient(client)
		want := client
		assertRegisteredClientId(t, got, want, err)
	})

	t.Run("remove client from registered list", func(t *testing.T) {
		hub := draw.NewHub()
		client := draw.NewClient(nil, hub)

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
	t.Run("received message from register channel", func(t *testing.T) {
		hub := draw.NewHub()
		client := draw.NewClient(nil, hub)

		go hub.Run()

		defer close(client.GetRegChan())
		client.SendRegistration()

		got, err := hub.GetClient(client)
		assertRegisteredClientId(t, got, client, err)

	})
	// t.Run("received message from un-register channel", func(t *testing.T) {
	// 	hub := draw.NewHub()
	// 	client := draw.NewClient(nil, hub)

	// 	// Run goroutine of hub listening for inputs from register/unresiger/broadcast
	// 	go hub.Run()

	// 	hub.RegisterClient(client)
	// 	got, err := hub.GetClient(client)
	// 	want := client
	// 	assertRegisteredClientId(t, got, want, err)

	// 	// Send unregistration and wait  for hub to get message from client
	// 	defer close(client.GetUnregChan())
	// 	client.SendUnregistration()

	// 	got, _ = hub.GetClient(client)
	// 	assertUnregisteredClient(t, got)
	// })

	t.Run("recieved payload from client to hub", func(t *testing.T) {
		hub := draw.NewHub()

		go hub.Run()

		c1 := draw.NewClient(nil, hub)
		c2 := draw.NewClient(nil, hub)
		c3 := draw.NewClient(nil, hub)
		c4 := draw.NewClient(nil, hub)

		c1.SendRegistration()
		c2.SendRegistration()
		c3.SendRegistration()
		c4.SendRegistration()

		c1.SendUpdate("Hello world1")

		assertPayloadSent(t, string(<-c1.GetSendChan()), "Hello world1")
		assertPayloadSent(t, string(<-c2.GetSendChan()), "Hello world1")
		assertPayloadSent(t, string(<-c3.GetSendChan()), "Hello world1")
		assertPayloadSent(t, string(<-c4.GetSendChan()), "Hello world1")

		c2.SendUpdate("Hello world2")

		assertPayloadSent(t, string(<-c1.GetSendChan()), "Hello world2")
		assertPayloadSent(t, string(<-c2.GetSendChan()), "Hello world2")
		assertPayloadSent(t, string(<-c3.GetSendChan()), "Hello world2")
		assertPayloadSent(t, string(<-c4.GetSendChan()), "Hello world2")

		c3.SendUpdate("Hello world3")

		assertPayloadSent(t, string(<-c1.GetSendChan()), "Hello world3")
		assertPayloadSent(t, string(<-c2.GetSendChan()), "Hello world3")
		assertPayloadSent(t, string(<-c3.GetSendChan()), "Hello world3")
		assertPayloadSent(t, string(<-c4.GetSendChan()), "Hello world3")

		c4.SendUpdate("Hello world4")

		assertPayloadSent(t, string(<-c1.GetSendChan()), "Hello world4")
		assertPayloadSent(t, string(<-c2.GetSendChan()), "Hello world4")
		assertPayloadSent(t, string(<-c3.GetSendChan()), "Hello world4")
		assertPayloadSent(t, string(<-c4.GetSendChan()), "Hello world4")

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

func assertPayloadSent(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("client did not recieve payload. got: %v, want: %v", got, want)
	}
}
