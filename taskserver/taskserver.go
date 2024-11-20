package taskserver

import (
	"encoding/json"
	"net/http"

	"github.com/manasm11/tasklist/taskstore"
)

func NewTaskServer() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]taskstore.Task{})
	})
}
