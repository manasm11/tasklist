package taskstore

type TaskStore struct{}

func New() *TaskStore {
	return &TaskStore{}
}
