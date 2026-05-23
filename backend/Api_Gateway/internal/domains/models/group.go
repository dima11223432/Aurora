package models

// Group represents a named collection of contacts.
type Group struct {
	ID       int64  `json:"id" db:"id"`
	Name     string `json:"group_name" db:"group_name"`
	Contacts []Contact
}
