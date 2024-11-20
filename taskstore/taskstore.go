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
	Id    uint64    `json:"id"`
	Title string    `json:"title"`
	Due   time.Time `json:"due"`
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
	task := Task{Id: idCounter, Title: title, Due: due}
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

func (ts *TaskStore) GetTasksBytDueDate(due time.Time) []Task {
	for _, task := range ts.tasks {
		if isDateEqual(due, task.Due) {
			return []Task{task}
		}
	}
	return nil
}

func (ts *TaskStore) DeleteTask(id uint64) error {
	_, ok := ts.tasks[id]
	if !ok {
		return TaskNotFoundError
	}
	delete(ts.tasks, id)
	for tag, ids := range ts.tagsTasks {
		for i, id := range ids {
			if id == id {
				ts.tagsTasks[tag] = append(ts.tagsTasks[tag][:i], ts.tagsTasks[tag][i+1:]...)
				break
			}
		}
	}
	return nil
}

func (ts *TaskStore) DeleteAllTasks() {
	ts.tasks = make(map[uint64]Task)
	ts.tagsTasks = make(map[string][]uint64)
}

func isDateEqual(t1, t2 time.Time) bool {
	return t1.Day() == t2.Day() && t1.Month() == t2.Month() && t1.Year() == t2.Year()
}
