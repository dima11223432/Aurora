package models

type Notification struct {
	ID      int64  `json:"id" db:"id"`
	Title   string `json:"notification_title" db:"notification_title"`
	Text    string `json:"notification_text"db:"notificaiton_text"`
	GroupID int64  `json:"group_id" db:"group_id"`
}
