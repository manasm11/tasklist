package taskserver

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/manasm11/tasklist/taskstore"
)

func Test_task(t *testing.T) {
	t.Run("get at task without any data", func(t *testing.T) {
		var ts http.Handler = NewTaskServer()
		h := httptest.NewServer(ts)
		defer h.Close()

		resp, err := http.Get(h.URL + "/task/")

		if err != nil {
			t.Errorf("got an error %v", err)
		}

		if resp.StatusCode != 200 {
			t.Errorf("expected 200 got %v", resp.StatusCode)
		}

		tasks := []taskstore.Task{}
		err = json.NewDecoder(resp.Body).Decode(&tasks)

		if err != nil {
			t.Errorf("got an error %v", err)
		}
	})

	t.Run("put at task without any data", func(t *testing.T) {
		var ts http.Handler = NewTaskServer()
		h := httptest.NewServer(ts)
		defer h.Close()

		req, err := http.NewRequest("PUT", h.URL+"/task/", nil)
		if err != nil {
			t.Errorf("got an error %v", err)
		}
		resp, err := http.DefaultClient.Do(req)

		if err != nil {
			t.Errorf("got an error %v", err)
		}
		if resp.StatusCode != http.StatusMethodNotAllowed {
			t.Errorf("expected %v got %v", http.StatusMethodNotAllowed, resp.StatusCode)
		}
	})
}


