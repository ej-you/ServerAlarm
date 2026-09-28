package entity

import "time"

// Climate describes server room climate data.
type Climate struct {
	ID          int
	Datetime    time.Time
	Temperature float32
}
