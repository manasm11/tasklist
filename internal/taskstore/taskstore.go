package taskstore

import "time"

type TaskStoreError string

func (terr TaskStoreError) Error() string {
	return string(terr)
}

const (
	TaskNotFoundError TaskStoreError = "task not found"
)

type Task struct {
	Id    uint64
	Title string
}

var idCounter uint64 = 0

type TaskStore struct {
	tasks     map[uint64]Task
	tagsTasks map[string][]uint64
}

func New() *TaskStore {
	tasks := make(map[uint64]Task)
	tagsTasks := make(map[string][]uint64)
	return &TaskStore{tasks, tagsTasks}
}

func (ts *TaskStore) CreateTask(title string, tags []string, due time.Time) uint64 {
	idCounter++
	task := Task{Id: idCounter, Title: title}
	ts.tasks[idCounter] = task
	for _, tag := range tags {
		ts.tagsTasks[tag] = append(ts.tagsTasks[tag], idCounter)
	}
	return task.Id
}

func (ts *TaskStore) GetTask(id uint64) (Task, error) {
	task, ok := ts.tasks[id]
	if !ok {
		return Task{}, TaskNotFoundError
	}
	return task, nil
}

func (ts *TaskStore) GetTasksByTag(tags string) (tasks []Task) {
	for _, id := range ts.tagsTasks[tags] {
		task := ts.tasks[id]
		tasks = append(tasks, task)
	}
	return tasks
}

func (ts *TaskStore) GetAllTask() []Task {
	tasks := make([]Task, len(ts.tasks))
	i := 0
	for _, task := range ts.tasks {
		tasks[i] = task
		i++
	}
	return tasks
}
