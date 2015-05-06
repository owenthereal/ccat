package db_common

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// NullInt represents an int that may be null. NullInt implements the Scanner
// interface so it can be used as a scan destination, similar to NullString.
type NullInt struct {
	Int   int
	Valid bool // Valid is true if Int is not NULL
}

// Scan implements the Scanner interface.
func (n *NullInt) Scan(value interface{}) error {
	if value == nil {
		n.Int, n.Valid = 0, false
		return nil
	}
	switch value := value.(type) {
	case int64:
		n.Int = int(value)
	default:
		return fmt.Errorf("scanning %T, got %T", n, value)
	}
	return nil
}

// Value implements the driver Valuer interface.
func (n NullInt) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Int, nil
}

// MarshalJSON implements the encoding/json.Marshaler interface.
func (i NullInt) MarshalJSON() ([]byte, error) {
	if i.Valid {
		return json.Marshal(i.Int)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON implements the encoding/json.Unmarshaler interface.
func (i *NullInt) UnmarshalJSON(data []byte) error {
	var v interface{}
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	if v == nil {
		*i = NullInt{}
	} else if v, ok := v.(float64); ok {
		*i = NullInt{Valid: true, Int: int(v)}
	} else {
		return fmt.Errorf("unmarshaling %T, got %T", i, v)
	}
	return nil
}

type NullTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

// Scan implements the Scanner interface.
func (nt *NullTime) Scan(value interface{}) error {
	nt.Time, nt.Valid = value.(time.Time)
	return nil
}

// Value implements the driver Valuer interface.
func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}

// MarshalJSON implements the json.Marshaler interface.
func (nt NullTime) MarshalJSON() ([]byte, error) {
	if nt.Valid {
		return json.Marshal(nt.Time)
	}
	return []byte("null"), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (nt *NullTime) UnmarshalJSON(data []byte) (err error) {
	if nt == nil {
		return errors.New("UnmarshalJSON on nil *NullTime pointer")
	}
	if bytes.Compare(data, []byte("null")) == 0 {
		nt.Valid = false
	} else {
		nt.Valid = true
		err = json.Unmarshal(data, &nt.Time)
	}
	return
}

func (nt NullTime) String() string {
	if nt.Valid {
		return nt.Time.String()
	}
	return "<nil>"
}

// Now returns a valid NullTime with the time set to now, in UTC and
// rounded to the nearest millisecond. It does this so that JSON- and
// SQL-serialization and deserialization yields the same time as passed
// in. If the time is not UTC or has sub-millisecond accuracy, the time
// retrieved from an SQL DB or JSON object might not be equal to the
// original object due to rounding and automatic timezone conversion.
func Now() NullTime {
	return NullTime{Time: time.Now().In(time.UTC).Round(time.Millisecond), Valid: true}
}

type StringSlice struct {
	Slice []string
}

func (s *StringSlice) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
	}

	inner := make([]string, len(s.Slice))
	for i, elem := range s.Slice {
		if strings.TrimSpace(elem) == "" || strings.Contains(elem, `"`) {
			inner[i] = strconv.Quote(elem)
		} else {
			inner[i] = elem
		}
	}
	return []byte("{" + strings.Join(inner, ",") + "}"), nil
}

func (s *StringSlice) Scan(v interface{}) error {
	if data, ok := v.([]byte); ok {
		interior := strings.Trim(string(data), "{}")
		if interior != "" {
			rawElems := strings.Split(interior, ",")
			s.Slice = make([]string, len(rawElems))
			for r, raw := range rawElems {
				if elem, err := strconv.Unquote(raw); err == nil {
					s.Slice[r] = elem
				} else {
					s.Slice[r] = raw
				}
			}
		} else {
			s.Slice = []string{}
		}
		return nil
	}
	return fmt.Errorf("%T.Scan failed: %v", s, v)
}

func NewSlice(goslice []string) *StringSlice {
	return &StringSlice{Slice: goslice}
}
