package racer

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRacer(t *testing.T) {
	t.Run("fastest server wins", func(t *testing.T) {
		slowServer := makeDelayedServerHelper(20 * time.Millisecond)
		fastServer := makeDelayedServerHelper(0 * time.Millisecond)

		defer slowServer.Close()
		defer fastServer.Close()

		slowURL := slowServer.URL
		fastURL := fastServer.URL

		want := fastURL
		got, _ := Racer(slowURL, fastURL)

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
	t.Run("return an error if server doesn't respond within 10s", func(t *testing.T) {
		server := makeDelayedServerHelper(25 * time.Millisecond)

		defer server.Close()

		_, err := ConfigurableRacer(server.URL, server.URL, 20*time.Millisecond)

		if err == nil {
			t.Errorf("expected an error but didn't get one")
		}
	})
}

func makeDelayedServerHelper(delay time.Duration) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		time.Sleep(delay)
		rw.WriteHeader(http.StatusOK)
	}))
}
