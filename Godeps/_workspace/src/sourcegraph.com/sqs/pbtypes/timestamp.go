package pbtypes

import "time"

// NewTimestamp creates a new Timestamp from a time.Time.
func NewTimestamp(t time.Time) Timestamp {
	return Timestamp{Seconds: t.Unix(), Nanos: int32(t.Nanosecond())}
}

func (t Timestamp) Time() time.Time {
	return time.Unix(t.Seconds, int64(t.Nanos))
}
