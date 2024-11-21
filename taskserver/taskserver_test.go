package taskserver

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/manasm11/tasklist/taskstore"
)

func Test_randompath(t *testing.T) {
	h := httptest.NewServer(NewTaskServer())
	defer h.Close()
	resp := reqWithoutData(t, h, "GET", "/randompath/")
	assertEqual(t, resp.StatusCode, http.StatusNotFound)
}

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

		getAllTasks(t, h)
		resp2 := reqWithoutData(t, h, "GET", "/task/")
		assertEqual(t, resp2.Header.Get("Content-Type"), "application/json")

		tasks := []taskstore.Task{}
		json.NewDecoder(resp2.Body).Decode(&tasks)

		assertEqual(t, len(tasks), 1)
		assertEqual(t, tasks[0].Title, "Task 1")
	})

	t.Run("delete at deletes all tasks", func(t *testing.T) {
		h := httptest.NewServer(NewTaskServer())
		defer h.Close()

		task1 := map[string]string{"title": "Task 1"}
		task2 := map[string]string{"title": "Task 2"}
		reqWithJsonData(t, h, "POST", "/task/", task1)
		reqWithJsonData(t, h, "POST", "/task/", task2)

		tasks := getAllTasks(t, h)
		assertEqual(t, len(tasks), 2)

		resp := reqWithoutData(t, h, "DELETE", "/task/")
		assertEqual(t, resp.StatusCode, http.StatusNoContent)

		tasks2 := getAllTasks(t, h)
		assertEqual(t, len(tasks2), 0)
	})
}

func Test_task_id(t *testing.T) {
	h := httptest.NewServer(NewTaskServer())
	defer h.Close()
	ids := createTasks(t, h, "Task 1", "Task 2", "Task 3")
	assertEqual(t, len(getAllTasks(t, h)), 3)
	for i, id := range ids {
		resp := reqWithoutData(t, h, "GET", "/task/"+strconv.FormatUint(id, 10)+"/")
		assertEqual(t, resp.StatusCode, http.StatusOK)
		task := taskstore.Task{}
		err := json.NewDecoder(resp.Body).Decode(&task)
		assertEqual(t, err, nil)
		assertEqual(t, task.Title, "Task "+strconv.Itoa(i+1))
	}
}

func createTasks(t testing.TB, h *httptest.Server, titles ...string) (ids []uint64) {
	t.Helper()
	for _, title := range titles {
		task := map[string]string{"title": title}
		resp := reqWithJsonData(t, h, "POST", "/task/", task)
		bs, err := io.ReadAll(resp.Body)
		assertEqual(t, err, nil)
		id, err := strconv.ParseInt(strings.Trim(string(bs), "\n"), 10, 64)
		assertEqual(t, err, nil)
		ids = append(ids, uint64(id))
	}
	return ids
}

func getAllTasks(t testing.TB, h *httptest.Server) []taskstore.Task {
	t.Helper()
	resp := reqWithoutData(t, h, "GET", "/task/")
	tasks := []taskstore.Task{}
	err := json.NewDecoder(resp.Body).Decode(&tasks)
	assertEqual(t, err, nil)
	return tasks
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
