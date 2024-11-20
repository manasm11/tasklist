package taskserver

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/manasm11/tasklist/taskstore"
)

type taskServer struct {
	ts *taskstore.TaskStore
}

func NewTaskServer() http.Handler {
	ts := taskstore.New()
	t := taskServer{ts: ts}
	//mux := http.NewServeMux()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/task/" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		switch r.Method {
		case "GET":
			t.taskGet(w, r)
		case "POST":
			t.taskPost(w, r)
		case "DELETE":
			t.taskDelete(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
}

func (t *taskServer) taskGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(t.ts.GetAllTask())
}

func (t *taskServer) taskPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	var task taskstore.Task
	json.NewDecoder(r.Body).Decode(&task)
	id := t.ts.CreateTask(task.Title, nil, time.Time{})
	json.NewEncoder(w).Encode(id)

}

func (t *taskServer) taskDelete(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
	t.ts.DeleteAllTasks()
}
