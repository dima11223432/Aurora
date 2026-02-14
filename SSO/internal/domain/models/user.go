package models

type User struct {
	ID       int
	Email    string
	IsAdmin  bool
	PassHash []byte
}
