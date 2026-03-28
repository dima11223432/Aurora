package services

import (
	"API_Service/internal/broker/kafka"
	"API_Service/internal/cache"
	"API_Service/internal/domains/models"
	"API_Service/internal/storage/postgres"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	// "github.com/sirupsen/logrus"
)

type ApiService struct {
	storage   *postgres.Storage
	publisher *kafka.Producer
	cache     *cache.RedisCache
}

func New(storage *postgres.Storage, publisher *kafka.Producer, redisCache *cache.RedisCache) *ApiService {
	return &ApiService{
		storage:   storage,
		publisher: publisher,
		cache:     redisCache,
	}
}

func (a *ApiService) SendNotification(ctx context.Context, notificationID int64) (string, error) {
	// notif, err := a.storage.GetNotificationByID(ctx, notificationID)
	// if err != nil {
	// 	return "", err
	// }
	var notif models.Notification
	err := a.cache.GetValue(ctx, fmt.Sprintf("%d", notificationID), &notif)
	if err != nil {
		if errors.Is(err, cache.ErrCacheMiss) {
			Dbnotif, err := a.storage.GetNotificationByID(ctx, notificationID)
			if err != nil {
				return "", err
			}
			notif = *Dbnotif
			go func(n models.Notification) {
				AsyncCtx := context.Background()
				a.cache.SetValue(AsyncCtx, fmt.Sprintf("%d", notificationID), n, time.Duration(10)*time.Minute)
			}(*Dbnotif)
		} else {
			return "", fmt.Errorf("cache error: %w", err)
		}
	}

	emails, err := a.storage.GetUsersEmailByGroupID(ctx, notif.GroupID)
	if err != nil {
		return "", err
	}

	event := models.SendNotificationEvent{
		EventID:        uuid.NewString(),
		EventType:      "send_notification",
		GroupID:        notif.GroupID,
		Emails:         emails,
		NotificationID: notificationID,
		CreatedAt:      time.Now(),
	}
	payload, err := json.Marshal(event)
	if err != nil {
		return "", err
	}

	a.publisher.Publish(ctx, "notification.sent", payload)
	return event.EventID, nil
}

func (a *ApiService) GetAllGroups(ctx context.Context) ([]models.Group, error) {
	return a.storage.GetAllGroups(ctx)
}

func (a *ApiService) UpdateGroup(ctx context.Context, groupID int64, newGroupName string) error {
	return a.storage.UpdateGroup(ctx, groupID, newGroupName)
}

func (a *ApiService) DeleteGroup(ctx context.Context, groupID int64) error {
	if groupID <= 0 {
		return fmt.Errorf("invalid groupID")
	}
	return a.storage.DeleteGroup(ctx, groupID)
}

func (a *ApiService) CreateGroup(ctx context.Context, groupName string) (int64, error) {
	id, err := a.storage.CreateGroup(ctx, groupName)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (a *ApiService) UpdateContact(ctx context.Context, contactID int64, newEmail string) error {
	return a.storage.UpdateContact(ctx, contactID, newEmail)
}

func (a *ApiService) DeleteContact(ctx context.Context, contactID int64) error {
	if contactID <= 0 {
		return fmt.Errorf("invalid contactID")
	}
	return a.storage.DeleteContact(ctx, contactID)
}

func (a *ApiService) GetAllContactsByGroupID(ctx context.Context, groupID int64) ([]models.Contact, error) {
	contacts, err := a.storage.GetContactsByGroupID(ctx, groupID)
	if err != nil {
		return nil, err
	}
	return contacts, nil
}

func (a *ApiService) CreateContact(
	ctx context.Context,
	email string,
	groupID int64,
) (int64, error) {

	return a.storage.CreateContact(ctx, models.Contact{
		Email:   email,
		GroupID: groupID,
	})
}

func (a *ApiService) CreateNotification(
	ctx context.Context,
	title string,
	text string,
	groupID int64,
) (int64, error) {

	return a.storage.CreateNotification(ctx, models.Notification{
		Title:   title,
		Text:    text,
		GroupID: groupID,
	})
}
