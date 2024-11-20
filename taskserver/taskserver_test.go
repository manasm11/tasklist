package taskserver

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/manasm11/tasklist/taskstore"
)

func Test_task(t *testing.T) {
	t.Run("get at task without any data", func(t *testing.T) {
		var ts http.Handler = NewTaskServer()
		h := httptest.NewServer(ts)
		defer h.Close()

		resp := reqWithoutData(t, h, "GET", "/task/")

		assertEqual(t, resp.StatusCode, 200)

		tasks := []taskstore.Task{}
		err := json.NewDecoder(resp.Body).Decode(&tasks)

		assertEqual(t, err, nil)
		assertEqual(t, len(tasks), 0)
	})

	t.Run("put at task without any data", func(t *testing.T) {
		h := httptest.NewServer(NewTaskServer())
		defer h.Close()

		resp := reqWithoutData(t, h, "PUT", "/task/")

		assertEqual(t, resp.StatusCode, http.StatusMethodNotAllowed)
	})

	t.Run("post at task create task", func(t *testing.T) {
		h := httptest.NewServer(NewTaskServer())
		defer h.Close()

		data := map[string]interface{}{"title": "Task 1"}
		resp := reqWithJsonData(t, h, "POST", "/task/", data)

		bs, err := io.ReadAll(resp.Body)
		assertEqual(t, err, nil)
		assertEqual(t, strings.Trim(string(bs), "\n"), "1")
		assertEqual(t, resp.StatusCode, http.StatusCreated)

		resp = reqWithoutData(t, h, "GET", "/task/")

		tasks := []taskstore.Task{}
		json.NewDecoder(resp.Body).Decode(&tasks)

		assertEqual(t, len(tasks), 1)
		assertEqual(t, tasks[0].Title, "Task 1")
	})
}

func assertEqual(t testing.TB, actual, expected interface{}) {
	t.Helper()
	if expected != actual {
		t.Errorf("expected %v got %v", expected, actual)
	}
}

func reqWithoutData(t testing.TB, h *httptest.Server, method, url string) *http.Response {
	t.Helper()
	req, err := http.NewRequest(method, h.URL+url, nil)
	if err != nil {
		t.Errorf("got an error %v", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Errorf("got an error %v", err)
	}
	return resp
}

func reqWithJsonData(t testing.TB, h *httptest.Server, method, url string, data interface{}) *http.Response {
	t.Helper()

	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Errorf("got an error %v", err)
	}

	req, err := http.NewRequest(method, h.URL+url, bytes.NewBuffer(jsonData))
	if err != nil {
		t.Errorf("got an error %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Errorf("got an error %v", err)
	}
	return resp
}
