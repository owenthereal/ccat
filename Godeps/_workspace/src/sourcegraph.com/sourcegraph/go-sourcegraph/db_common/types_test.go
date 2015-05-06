package db_common

import (
	"encoding/json"
	"testing"
	"time"
)

func TestNullTime(t *testing.T) {
	tests := []struct {
		t    NullTime
		json string
	}{
		{NullTime{Valid: false}, "null"},
		{NullTime{Time: time.Unix(12345, 0).In(time.UTC), Valid: true}, `"1970-01-01T03:25:45Z"`},
	}
	for _, test := range tests {
		tjson, err := json.Marshal(test.t)
		if err != nil {
			t.Errorf("NullTime %+v: failed to marshal: %s", test.t, err)
		}
		if test.json != string(tjson) {
			t.Errorf("NullTime %+v: want JSON %q, got %q", test.t, test.json, string(tjson))
		}
		var tp NullTime
		err = json.Unmarshal(tjson, &tp)
		if err != nil {
			t.Errorf("NullTime %+v: failed to unmarshal: %s", test.t, err)
		}
		if test.t != tp {
			t.Errorf("NullTime %+v: marshal-then-unmarshal produced different obj %+v", test.t, tp)
		}
	}
}
