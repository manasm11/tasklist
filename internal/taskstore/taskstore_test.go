package taskstore

import "testing"

func TestNew(t *testing.T) {
	var ts *TaskStore = New()

	if ts == nil {
		t.Errorf("New() returned nil")
	}
}
