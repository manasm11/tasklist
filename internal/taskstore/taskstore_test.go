package taskstore

import (
	"log"
	"testing"
)

func TestNew(t *testing.T) {
	var ts *TaskStore = New()

	if ts == nil {
		t.Errorf("New() returned nil")
	}
}

func TestCreateTask(t *testing.T) {

	var ts *TaskStore = New()

	var id uint64 = ts.CreateTask("Task 1")

	task, err := ts.GetTask(id)
	if err != nil {
		t.Errorf("GetTask() returned error: %v", err)
	}

	if task.Id != id {
		t.Errorf("incorrect task id: %v", task.Id)
	}

	if task.Title != "Task 1" {
		t.Errorf("incorrect task title: %v", task.Title)
	}
}

func TestGetTask(t *testing.T) {
	t.Run("add and get 1 task", func(t *testing.T) {
		var ts *TaskStore = New()
		id := ts.CreateTask("Task 1")

		task, err := ts.GetTask(id)

		if err != nil {
			t.Errorf("GetTask() returned error: %v", err)
		}

		if task.Title != "Task 1" {
			t.Errorf("incorrect task title: %v", task.Title)
		}
	})

	t.Run("add and get 2 tasks", func(t *testing.T) {
		var ts *TaskStore = New()
		id1 := ts.CreateTask("Task 1")
		id2 := ts.CreateTask("Task 2")

		task1, err := ts.GetTask(id1)

		if err != nil {
			t.Errorf("GetTask() returned error: %v", err)
		}

		if task1.Id != id1 {
			t.Errorf("incorrect task id: %v", task1.Id)
		}

		if task1.Title != "Task 1" {
			t.Errorf("incorrect task title: %v", task1.Title)
		}

		task2, err := ts.GetTask(id2)

		if err != nil {
			t.Errorf("GetTask() returned error: %v", err)
		}

		if task2.Id != id2 {
			t.Errorf("incorrect task id: %v", task2.Id)
		}

		if task2.Title != "Task 2" {
			t.Errorf("incorrect task title: %v", task2.Title)
		}
	})
}

func TestGetAllTask(t *testing.T) {
	var ts *TaskStore = New()
	ts.CreateTask("Task 1")

	var tasks []Task = ts.GetAllTask()

	log.Println("ts.tasks:", ts.tasks)
	log.Println("tasks:", tasks)

	if len(tasks) != 1 {
		t.Errorf("GetAllTask() returned %v tasks", len(tasks))
	}

	if tasks[0].Title != "Task 1" {
		t.Errorf("incorrect task title: %v", tasks[0].Title)
	}
}
