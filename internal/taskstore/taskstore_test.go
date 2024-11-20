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

	var id uint64 = ts.CreateTask("Task 1", nil)

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
		id := ts.CreateTask("Task 1", nil)

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
		id1 := ts.CreateTask("Task 1", nil)
		id2 := ts.CreateTask("Task 2", nil)

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
	ts.CreateTask("Task 1", nil)

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

func TestGetTasksByTag(t *testing.T) {
	t.Run("create and access 1 task with tag", func(t *testing.T) {
		ts := New()
		id := ts.CreateTask("Task 1", []string{"tag1"})

		tasks := ts.GetTasksByTag("tag1")

		if tasks[0].Title != "Task 1" {
			t.Errorf("incorrect task title: %v", tasks[0].Title)
		}

		if tasks[0].Id != id {
			t.Errorf("incorrect task id: %v", tasks[0].Id)
		}
	})

	t.Run("create two tasks with different tags and access just one", func(t *testing.T) {
		ts := New()
		id1 := ts.CreateTask("Task 1", []string{"tag1"})
		id2 := ts.CreateTask("Task 1", []string{"tag1", "tag2"})

		tag1Tasks := ts.GetTasksByTag("tag1")
		tag2Tasks := ts.GetTasksByTag("tag2")

		if len(tag1Tasks) != 2 {
			t.Errorf("GetTasksByTag() returned %v tasks", len(tag1Tasks))
		}

		if len(tag2Tasks) != 1 {
			t.Errorf("GetTasksByTag() returned %v tasks", len(tag2Tasks))
		}

		if tag2Tasks[0].Id != id2 {
			t.Errorf("incorrect task id: %v", tag2Tasks[0].Id)
		}

		if tag1Tasks[0].Id != id1 && tag1Tasks[1].Id != id1 {
			t.Errorf("task id %d not in %v", id1, tag1Tasks)
		}

		if tag1Tasks[0].Id != id2 && tag1Tasks[1].Id != id2 {
			t.Errorf("task id %d not in %v", id2, tag1Tasks)
		}
	})
}
