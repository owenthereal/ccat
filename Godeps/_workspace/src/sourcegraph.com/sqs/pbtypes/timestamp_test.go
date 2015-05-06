package pbtypes

import (
	"testing"
	"time"
)

func TestTimestamp(t *testing.T) {
	tm := time.Now()
	ts := NewTimestamp(tm)
	tm2 := ts.Time()
	if !tm2.Equal(tm) {
		t.Errorf("got %q, want %q", tm2, tm)
	}
}
