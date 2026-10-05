package threshold

import "cmp"

const (
	Zero = iota
	Min
	Low
	Standart
	High
	Urgent
)

// Threshold represents an util helps to determine threshold level.
type Threshold[T cmp.Ordered] struct {
	standart T
	high     T
	urgent   T
}

// New returns a new instance of Threshold.
func New[T cmp.Ordered](standart, high, urgent T) *Threshold[T] {
	return &Threshold[T]{
		standart: standart,
		high:     high,
		urgent:   urgent,
	}
}

// Level returns threshold level for the given value.
func (t *Threshold[T]) Level(value T) int {
	switch {
	case value >= t.urgent:
		return Urgent
	case value >= t.high:
		return High
	case value >= t.standart:
		return Standart
	default:
		return Zero
	}
}

// LevelString returns string representation of theshold level.
func LevelString(level int) string {
	switch level {
	case Urgent:
		return "urgent"
	case High:
		return "high"
	case Standart:
		return "default"
	case Low:
		return "low"
	case Min:
		return "min"
	default:
		return "no"
	}
}
