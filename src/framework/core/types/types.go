package types

import (
	"time"
)

// MapStr the common event data definition
type MapStr map[string]interface{}

// Event the event data definition
type Event struct {
	Timestamp time.Time `json:"timestamp"`
	Fields    MapStr    `json:",inline"`
}
