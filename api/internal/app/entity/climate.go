package entity

import (
	"fmt"
	"time"
)

// Climate describes server room climate data.
type Climate struct {
	ID          int
	Datetime    time.Time
	Temperature float32
}

// TemperatureString returns string representation of climate temperature.
func (c *Climate) TemperatureString() string {
	return fmt.Sprintf("%.1f°C", c.Temperature)
}
