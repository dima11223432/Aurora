package models

// App represents an application that can authenticate users.
type App struct {
	ID     int
	Name   string
	Secret string
}
