package models

import (
	"time"
)

type GroupCreateEvent struct {
	EventID   string    `json:"event_id"`
	GroupID   int64     `json:"group_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type SendNotificationEvent struct {
	EventID        string    `json:"event_id"`
	EventType      string    `json:"event_type"`
	GroupID        int64     `json:"group_id"`
	Emails         []string  `json:"emails"`
	NotificationID int64     `json:"notification_id"`
	CreatedAt      time.Time `json:"created_at"`
}
