package types

// MapStr the common event data definition
type MapStr map[string]interface{}

// Saver the save interface
type Saver interface {
	Save() error
}
