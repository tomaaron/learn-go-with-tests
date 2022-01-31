package context

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type StubStore struct {
	response  string
	cancelled bool
}

func (s *StubStore) Fetch() string {
	time.Sleep(100 * time.Millisecond)
	return s.response
}

func (s *StubStore) Cancel() {
	s.cancelled = true
}
func TestServer(t *testing.T) {
	t.Run("returns data from store", func(t *testing.T) {
		data := "hello, world"
		store := &StubStore{response: data}
		svr := Server(store)

		request := httptest.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()

		svr.ServeHTTP(response, request)

		if response.Body.String() != data {
			t.Errorf(`got "%s", want "%s"`, response.Body.String(), data)
		}

		if store.cancelled {
			t.Error("it should not have cancelled the store")
		}
	})

	t.Run("tells store to cancel work if request is cancelled", func(t *testing.T) {
		data := "hello, world"
		store := &StubStore{response: data, cancelled: false}
		svr := Server(store)

		request := httptest.NewRequest(http.MethodGet, "/", nil)

		cancellingCtx, cancel := context.WithCancel(request.Context())
		time.AfterFunc(5*time.Millisecond, cancel)
		request = request.WithContext(cancellingCtx)

		response := httptest.NewRecorder()

		svr.ServeHTTP(response, request)

		if !store.cancelled {
			t.Error("store was not told to cancel")
		}
	})
}
