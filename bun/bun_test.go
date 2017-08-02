package bun

import (
	"testing"
	"errors"
)

type ErrorWriter struct{}

func (e *ErrorWriter) Write(b []byte) (int, error) {
	return 0, errors.New("Expected error")
}

func TestLogBun(t *testing.T) {
	var bun1, bun2 Bun

	bun1.LogBun(1, "")

	if bun1.Location == nil {
		t.Errorf("Bun has no location")
	}
	if bun1.Size == 0 {
		t.Errorf("Bun has 0 size")
	}
	if bun1.Description != "" {
		t.Errorf("Expected empty description, got %s", bun1.Description)
	}

	bun2.LogBun(2, "description")

	if bun2.Size != 2 {
		t.Errorf("Bun has size %d, should have size 2", bun2.Size)
	}
	if bun2.Description != "description" {
		t.Errorf("Expected description, got %s", bun2.Description)
	}
}
