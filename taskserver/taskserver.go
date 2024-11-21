package taskserver

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/manasm11/tasklist/taskstore"
)

type taskServer struct {
	ts *taskstore.TaskStore
}

func NewTaskServer() http.Handler {
	ts := taskstore.New()
	t := taskServer{ts: ts}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /task/", t.taskGet)
	mux.HandleFunc("POST /task/", t.taskPost)
	mux.HandleFunc("DELETE /task/", t.taskDelete)
	mux.HandleFunc("GET /task/{id}/", t.taskIdGet)
	return mux
}

func (t *taskServer) taskIdGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idString := r.PathValue("id")
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		http.Error(w, "id must be string", http.StatusBadRequest)
		return
	}
	uid := uint64(id)
	task, err := t.ts.GetTask(uid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
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
