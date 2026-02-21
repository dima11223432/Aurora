package models

type Group struct {
	ID       int64  `json:"id" db:"id"`
	Name     string `json:"group_name" db:"group_name"`
	Contacts []Contact
}
