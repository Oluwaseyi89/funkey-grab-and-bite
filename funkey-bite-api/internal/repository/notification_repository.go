package repository

import (
	"database/sql"
	"fmt"
	"time"

	"funkey-grab-and-bite/funkey-bite-api/internal/domain/models"
)

type NotificationRepository struct {
	db *sql.DB
}

func NewNotificationRepository(db *sql.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) Create(notification *models.Notification) error {
	query := `
		INSERT INTO notifications (user_id, type, title, message, is_read, reference_id, reference_type, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at
	`

	err := r.db.QueryRow(
		query,
		notification.UserID,
		notification.Type,
		notification.Title,
		notification.Message,
		notification.IsRead,
		notification.ReferenceID,
		notification.ReferenceType,
		time.Now(),
	).Scan(&notification.ID, &notification.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
	}

	return nil
}

func (r *NotificationRepository) GetByUserID(userID int, limit int) ([]models.Notification, error) {
	query := `
		SELECT id, user_id, type, title, message, is_read, reference_id, reference_type, created_at, read_at
		FROM notifications
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`

	rows, err := r.db.Query(query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get notifications: %w", err)
	}
	defer rows.Close()

	var notifications []models.Notification
	for rows.Next() {
		var notification models.Notification
		var referenceID sql.NullInt64
		var readAt sql.NullTime

		err := rows.Scan(
			&notification.ID,
			&notification.UserID,
			&notification.Type,
			&notification.Title,
			&notification.Message,
			&notification.IsRead,
			&referenceID,
			&notification.ReferenceType,
			&notification.CreatedAt,
			&readAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan notification: %w", err)
		}

		if referenceID.Valid {
			val := int(referenceID.Int64)
			notification.ReferenceID = &val
		}
		if readAt.Valid {
			notification.ReadAt = &readAt.Time
		}

		notifications = append(notifications, notification)
	}

	return notifications, nil
}

func (r *NotificationRepository) MarkAsRead(notificationID int) error {
	query := `UPDATE notifications SET is_read = true, read_at = $1 WHERE id = $2`
	_, err := r.db.Exec(query, time.Now(), notificationID)
	return err
}

func (r *NotificationRepository) GetUnreadCount(userID int) (int, error) {
	query := `SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND is_read = false`
	var count int
	err := r.db.QueryRow(query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get unread count: %w", err)
	}
	return count, nil
}

func (r *NotificationRepository) DeleteOldNotifications(days int) error {
	query := `DELETE FROM notifications WHERE created_at < NOW() - INTERVAL '$1 days'`
	_, err := r.db.Exec(query, days)
	return err
}
