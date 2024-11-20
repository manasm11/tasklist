package taskstore

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	var ts *TaskStore = New()

	if ts == nil {
		t.Errorf("New() returned nil")
	}
}

func TestCreateTask(t *testing.T) {

	var ts *TaskStore = New()

	var id uint64 = ts.CreateTask("Task 1", nil, time.Time{})

	task, err := ts.GetTask(id)
	if err != nil {
		t.Errorf("GetTask() returned error: %v", err)
	}

	assertTaskTitleId(t, task, "Task 1", id)
}

func TestGetTask(t *testing.T) {
	t.Run("add and get 1 task", func(t *testing.T) {
		ts, id := createTaskStoreAndTaskWithoutTagsOrDue("Task 1")

		task, err := ts.GetTask(id)

		if err != nil {
			t.Errorf("GetTask() returned error: %v", err)
		}

		assertTaskTitleId(t, task, "Task 1", id)
	})

	t.Run("add and get 2 tasks", func(t *testing.T) {
		ts, id1 := createTaskStoreAndTaskWithoutTagsOrDue("Task 1")
		id2 := ts.CreateTask("Task 2", nil, time.Time{})

		task1, err := ts.GetTask(id1)

		if err != nil {
			t.Errorf("GetTask() returned error: %v", err)
		}
		assertTaskTitleId(t, task1, "Task 1", id1)

		task2, err := ts.GetTask(id2)

		if err != nil {
			t.Errorf("GetTask() returned error: %v", err)
		}

		assertTaskTitleId(t, task2, "Task 2", id2)
	})
}

func TestGetAllTask(t *testing.T) {
	ts, id := createTaskStoreAndTaskWithoutTagsOrDue("Task 1")

	var tasks []Task = ts.GetAllTask()

	assertLen(t, tasks, 1)
	assertTaskTitleId(t, tasks[0], "Task 1", id)
}

func TestGetTasksByTag(t *testing.T) {
	t.Run("create and access 1 task without tag", func(t *testing.T) {
		ts, _ := createTaskStoreAndTaskWithoutTagsOrDue("Task 1")

		tasks := ts.GetTasksByTag("tag1")

		assertLen(t, tasks, 0)
	})

	t.Run("create and access 1 task with tag", func(t *testing.T) {
		ts := New()
		id := ts.CreateTask("Task 1", []string{"tag1"}, time.Time{})

		tasks := ts.GetTasksByTag("tag1")
		assertLen(t, tasks, 1)
		assertTaskTitleId(t, tasks[0], "Task 1", id)
	})

	t.Run("create two tasks with different tags and access just one", func(t *testing.T) {
		ts := New()
		id1 := ts.CreateTask("Task 1", []string{"tag1"}, time.Time{})
		id2 := ts.CreateTask("Task 2", []string{"tag1", "tag2"}, time.Time{})

		tag1Tasks := ts.GetTasksByTag("tag1")
		tag2Tasks := ts.GetTasksByTag("tag2")

		assertLen(t, tag1Tasks, 2)
		assertLen(t, tag2Tasks, 1)

		assertTaskTitleId(t, tag2Tasks[0], "Task 2", id2)

		if tag1Tasks[0].Id != id1 && tag1Tasks[1].Id != id1 {
			t.Errorf("task id %d not in %v", id1, tag1Tasks)
		}

		if tag1Tasks[0].Id != id2 && tag1Tasks[1].Id != id2 {
			t.Errorf("task id %d not in %v", id2, tag1Tasks)
		}
	})
}

func TestGetTasksByDueDate(t *testing.T) {
	t.Run("create and access 1 task without due date", func(t *testing.T) {
		ts, _ := createTaskStoreAndTaskWithoutTagsOrDue("Task 1")
		tasks := ts.GetTasksBytDueDate(time.Now())

		assertLen(t, tasks, 0)
	})

	t.Run("create and access 1 task with due date", func(t *testing.T) {
		ts := New()
		id := ts.CreateTask("Task 1", nil, time.Now())
		time.Sleep(20 * time.Millisecond)
		tasks := ts.GetTasksBytDueDate(time.Now())

		assertLen(t, tasks, 1)
		assertTaskTitleId(t, tasks[0], "Task 1", id)
	})

	t.Run("create 2 tasks and access 1 task with due date", func(t *testing.T) {
		ts := New()
		today := time.Now()
		tomorrow := today.AddDate(0, 0, 1)
		id1 := ts.CreateTask("Task 1", nil, today)
		id2 := ts.CreateTask("Task 2", nil, tomorrow)
		tasks := ts.GetTasksBytDueDate(tomorrow)

		assertLen(t, tasks, 1)
		assertTaskTitleId(t, tasks[0], "Task 2", id2)

		tasks = ts.GetTasksBytDueDate(today)

		assertLen(t, tasks, 1)
		assertTaskTitleId(t, tasks[0], "Task 1", id1)
	})
}

func TestDeleteTask(t *testing.T) {
	t.Run("delete task that does not exist", func(t *testing.T) {
		ts := New()
		err := ts.DeleteTask(1)
		if err != TaskNotFoundError {
			t.Errorf("DeleteTask() returned wrong error: %v", err)
		}
	})

	t.Run("delete a task that exists", func(t *testing.T) {
		ts, id := createTaskStoreAndTaskWithoutTagsOrDue("Task 1")

		err := ts.DeleteTask(id)

		assertEqual(t, err, nil)
		assertLen(t, ts.GetAllTask(), 0)
	})

	t.Run("delete task from many tasks", func(t *testing.T) {
		ts, id1 := createTaskStoreAndTaskWithoutTagsOrDue("Task 1")
		_ = ts.CreateTask("Task 2", nil, time.Time{})

		err := ts.DeleteTask(id1)

		assertEqual(t, err, nil)
		assertLen(t, ts.GetAllTask(), 1)
	})
}

func assertEqual[T comparable](t testing.TB, got, want T) {
	t.Helper()
	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func assertLen(t testing.TB, slice []Task, expectedLen int) {
	t.Helper()
	if len(slice) != expectedLen {
		t.Errorf("got length %d want %d", len(slice), expectedLen)
	}
}

func assertTaskTitleId(t testing.TB, task Task, title string, id uint64) {
	t.Helper()

	if task.Id != id {
		t.Errorf("got id %d want %d", task.Id, id)
	}

	if task.Title != title {
		t.Errorf("got title %q want %q", task.Title, title)
	}
}

func createTaskStoreAndTaskWithoutTagsOrDue(title string) (ts *TaskStore, id uint64) {
	ts = New()
	id = ts.CreateTask(title, nil, time.Time{})
	return ts, id
}
