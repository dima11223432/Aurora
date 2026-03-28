package postgres

import (
	"API_Service/internal/domains/models"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "internal.storage.postgres.new"

	db, err := sql.Open("postgres", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Storage{
		db: db,
	}, nil
}

func (s *Storage) CreateGroup(ctx context.Context, groupName string) (int64, error) {
	const op = "internal.storage.postgres.CreateGroup"
	query := `
	INSERT INTO notification_groups (name) VALUES ($1) returning id
	`
	var id int64

	err := s.db.QueryRowContext(ctx, query, groupName).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil

}

func (s *Storage) DeleteGroup(ctx context.Context, groupID int64) error {
	const op = "internal.storage.postgres.DeleteGroup"

	query := `
	DELETE FROM notification_groups WHERE id=$1 returning id
	`

	var DeletedGroupID int64

	err := s.db.QueryRowContext(ctx, query, groupID).Scan(&DeletedGroupID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf(":%s group not found ", op)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetAllGroups(ctx context.Context) ([]models.Group, error) {
	const op = "internal.storage.postgres.GetAllGroups"
	query := `
	SELECT id, name from notification_groups
	`

	rows, err := s.db.QueryContext(ctx, query)
	defer rows.Close()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var groups []models.Group

	for rows.Next() {
		var id int64
		var name string
		err := rows.Scan(&id, &name)
		if err != nil {
			return nil, err
		}
		groups = append(groups, models.Group{
			ID:   id,
			Name: name,
		})
	}
	return groups, nil
}

func (s *Storage) UpdateGroup(ctx context.Context, groupID int64, newGroupName string) error {
	const op = "internal.storage.postgres.Updategroup"

	query := `
	UPDATE notification_groups SET name = $1 WHERE id = $2 RETURNING id
	`

	var UpdatedGroupID int64

	err := s.db.QueryRowContext(ctx, query, newGroupName, groupID).Scan(&UpdatedGroupID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("invalid groupID ")
		}
		return err
	}

	return nil
}

func (s *Storage) DeleteContact(ctx context.Context, contactID int64) error {
	const op = "internal.storage.postgres.DeleteContact"

	if contactID == 0 {
		return fmt.Errorf("invalid contactID")
	}

	query := `
	DELETE FROM contacts WHERE id=$1 RETURNING id
	`

	var DeletedContactID int64

	err := s.db.QueryRowContext(ctx, query, contactID).Scan(&DeletedContactID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("contact not found")
		}
		return err
	}
	return nil
}

func (s *Storage) UpdateContact(ctx context.Context, contactID int64, newEmail string) error {
	const op = "internal.storage.postgres.UpdateContact"

	query := `
	UPDATE contacts SET email = $1 WHERE id = $2 RETURNING id
	`

	var UpdatedContactID int64

	err := s.db.QueryRowContext(ctx, query, newEmail, contactID).Scan(&UpdatedContactID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("contact not found")
		}
		return err
	}
	return nil
}
func (s *Storage) CreateContact(ctx context.Context, contact models.Contact) (int64, error) {
	const op = "internal.storage.postgres.CreateContact"

	query := `
	INSERT INTO contacts (email, group_id)
	VALUES ($1, $2)
	RETURNING id
	`

	var id int64
	err := s.db.QueryRowContext(ctx, query, contact.Email, contact.GroupID).Scan(&id)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code.Name() == "foreign_key_violation" {
				return 0, fmt.Errorf("%s: group with id %d does not exists", op, contact.GroupID)
			}
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil

}

func (s *Storage) GetContactsByGroupID(ctx context.Context, groupID int64) ([]models.Contact, error) {
	const op = "internal.storage.postgres.GetContactsByGroupID"

	query := `
	SELECT id,email
	FROM contacts
	WHERE group_id = $1
	`

	rows, err := s.db.QueryContext(ctx, query, groupID)
	defer rows.Close()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	contacts := make([]models.Contact, 0)
	for rows.Next() {
		var email string
		var contactID int64
		if err := rows.Scan(&contactID, &email); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		contacts = append(contacts, models.Contact{
			ID:    contactID,
			Email: email,
		})
	}
	return contacts, nil
}

// TODO: написать функцию для создания массива email всех юзеров в группе с id == groupID
func (s *Storage) GetUsersEmailByGroupID(ctx context.Context, groupID int64) ([]string, error) {
	const op = "internal.storage.postgres.GetUsersEmailByGroupID"

	query := `
	SELECT email
	FROM contacts
	WHERE group_id = $1
	`

	rows, err := s.db.QueryContext(ctx, query, groupID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	emails := make([]string, 0)
	for rows.Next() {
		var email string
		if err := rows.Scan(&email); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		emails = append(emails, email)
	}
	return emails, nil
}

func (s *Storage) GetNotificationByID(ctx context.Context, notificationID int64) (*models.Notification, error) {
	const op = "internal.storage.postgres.GetNotificationByID"

	query := `
	SELECT id,title, text, group_id FROM notifications
	WHERE id = $1
	`

	var notification models.Notification

	err := s.db.QueryRowContext(ctx, query, notificationID).Scan(&notificationID, &notification.Title, &notification.Text, &notification.GroupID)
	if err != nil {
		return &models.Notification{}, fmt.Errorf("%s:%w", op, err)
	}

	return &notification, nil
}

func (s *Storage) CreateNotification(ctx context.Context, notification models.Notification) (int64, error) {
	const op = "internal.storage.postgres.CreateNotification"

	query := `
	INSERT INTO notifications (title, text, group_id)
	VALUES ($1, $2, $3)
	RETURNING id
	`

	var id int64
	err := s.db.QueryRowContext(ctx, query, notification.Title, notification.Text, notification.GroupID).Scan(&id)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code.Name() == "foreign_key_violation" {
				return 0, fmt.Errorf("%s: group with this id %d does not exists", op, notification.GroupID)
			}
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}
