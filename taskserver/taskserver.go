package taskserver

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/manasm11/tasklist/taskstore"
)

func NewTaskServer() http.Handler {
	ts := taskstore.New()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(ts.GetAllTask())
		case "POST":
			w.WriteHeader(http.StatusCreated)
			var t taskstore.Task
			json.NewDecoder(r.Body).Decode(&t)
			id := ts.CreateTask(t.Title, nil, time.Time{})
			json.NewEncoder(w).Encode(id)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
}
