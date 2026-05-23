package models

// Contact represents a contact with email associated to a group.
type Contact struct {
	ID      int64  `json:"id" db:"id"`
	Email   string `json:"email" db:"email"`
	GroupID int64  `json:"group_id" db:"group_id"`
}
