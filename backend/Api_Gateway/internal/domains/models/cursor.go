package models

// Cursor implements cursor-based pagination with a score and unique ID.
type Cursor struct {
	Score float64
	ID    string
}
