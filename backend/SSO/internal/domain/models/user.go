package models

type User struct {
	ID          int64
	Telegram_id int64
	Username    string
	First_name  string
	Last_name   string
	Is_admin    bool
}
