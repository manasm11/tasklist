package taskstore

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
	tasks map[uint64]Task
}

func New() *TaskStore {
	tasks := make(map[uint64]Task)
	return &TaskStore{tasks}
}

func (ts *TaskStore) CreateTask(title string) uint64 {
	idCounter++
	task := Task{Id: idCounter, Title: title}
	ts.tasks[idCounter] = task
	return task.Id
}

func (ts *TaskStore) GetTask(id uint64) (Task, error) {
	task, ok := ts.tasks[id]
	if !ok {
		return Task{}, TaskNotFoundError
	}
	return task, nil
}

func (ts *TaskStore) GetAllTask() []Task {
	tasks := make([]Task, len(ts.tasks))
	for _, task := range ts.tasks {
		tasks = append(tasks, task)
	}
	return tasks
}
