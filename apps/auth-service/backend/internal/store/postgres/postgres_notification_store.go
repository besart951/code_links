package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (s *Store) GetSmtpSettings(ctx context.Context) (SmtpSettings, error) {
	var settings SmtpSettings
	err := s.pool.QueryRow(ctx, `
		select host, port, username, coalesce(password_encrypted, ''), encryption, from_email, from_name, reply_to_email, active, updated_at
		from smtp_settings
		where id = 'default'
	`).Scan(&settings.Host, &settings.Port, &settings.Username, &settings.PasswordEncrypted, &settings.Encryption, &settings.FromEmail, &settings.FromName, &settings.ReplyToEmail, &settings.Active, &settings.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		settings.Port = 587
		settings.Encryption = SmtpEncryptionStartTLS
		settings.FromEmail = "no-reply@codelinks.dev"
		settings.FromName = "CodeLinks"
		settings.ReplyToEmail = "support@codelinks.dev"
		settings.UpdatedAt = time.Now().UTC()
		return settings, nil
	}
	settings.HasPassword = settings.PasswordEncrypted != ""

	return settings, err
}

func (s *Store) SaveSmtpSettings(ctx context.Context, settings SmtpSettings) (SmtpSettings, error) {
	current, _ := s.GetSmtpSettings(ctx)
	if settings.PasswordEncrypted == "" {
		settings.PasswordEncrypted = current.PasswordEncrypted
	}
	err := s.pool.QueryRow(ctx, `
		insert into smtp_settings (id, host, port, username, password_encrypted, encryption, from_email, from_name, reply_to_email, active, updated_at)
		values ('default', $1, $2, $3, $4, $5, $6, $7, $8, $9, now())
		on conflict (id) do update set
			host = excluded.host,
			port = excluded.port,
			username = excluded.username,
			password_encrypted = excluded.password_encrypted,
			encryption = excluded.encryption,
			from_email = excluded.from_email,
			from_name = excluded.from_name,
			reply_to_email = excluded.reply_to_email,
			active = excluded.active,
			updated_at = now()
		returning host, port, username, coalesce(password_encrypted, ''), encryption, from_email, from_name, reply_to_email, active, updated_at
	`, settings.Host, settings.Port, settings.Username, settings.PasswordEncrypted, settings.Encryption, settings.FromEmail, settings.FromName, settings.ReplyToEmail, settings.Active).
		Scan(&settings.Host, &settings.Port, &settings.Username, &settings.PasswordEncrypted, &settings.Encryption, &settings.FromEmail, &settings.FromName, &settings.ReplyToEmail, &settings.Active, &settings.UpdatedAt)
	settings.HasPassword = settings.PasswordEncrypted != ""

	return settings, err
}

func (s *Store) ListNotifications(ctx context.Context, limit int) ([]Notification, error) {
	rows, err := s.pool.Query(ctx, `
		select id, type, channel, recipient, subject, status, created_at, sent_at
		from notifications
		order by created_at desc
		limit $1
	`, normalizedLimit(limit))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notifications := []Notification{}
	for rows.Next() {
		var notification Notification
		if err := rows.Scan(&notification.ID, &notification.Type, &notification.Channel, &notification.Recipient, &notification.Subject, &notification.Status, &notification.CreatedAt, &notification.SentAt); err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)
	}

	return notifications, rows.Err()
}

func (s *Store) CreateNotification(ctx context.Context, notification Notification) error {
	if notification.ID == uuid.Nil {
		notification.ID = uuid.New()
	}
	if notification.CreatedAt.IsZero() {
		notification.CreatedAt = time.Now().UTC()
	}
	_, err := s.pool.Exec(ctx, `
		insert into notifications (id, type, channel, recipient, subject, status, created_at, sent_at)
		values ($1, $2, $3, $4, $5, $6, $7, $8)
	`, notification.ID, notification.Type, notification.Channel, notification.Recipient, notification.Subject, notification.Status, notification.CreatedAt, notification.SentAt)

	return err
}

func (s *Store) RecordAdminAuditEntry(ctx context.Context, entry AdminAuditEntry) error {
	if entry.ID == uuid.Nil {
		entry.ID = uuid.New()
	}
	if entry.CreatedAt.IsZero() {
		entry.CreatedAt = time.Now().UTC()
	}
	_, err := s.pool.Exec(ctx, `
		insert into admin_audit_entries (id, actor_user_id, action, target_type, target_id, reason, ip_address, created_at)
		values ($1, $2, $3, $4, $5, $6, nullif($7, '')::inet, $8)
	`, entry.ID, entry.ActorUserID, entry.Action, entry.TargetType, entry.TargetID, entry.Reason, entry.IPAddress, entry.CreatedAt)

	return err
}

func (s *Store) ListAdminAuditEntries(ctx context.Context, limit int) ([]AdminAuditEntry, error) {
	rows, err := s.pool.Query(ctx, `
		select id, actor_user_id, action, target_type, target_id, coalesce(reason, ''), coalesce(ip_address::text, ''), created_at
		from admin_audit_entries
		order by created_at desc
		limit $1
	`, normalizedLimit(limit))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entries := []AdminAuditEntry{}
	for rows.Next() {
		var entry AdminAuditEntry
		if err := rows.Scan(&entry.ID, &entry.ActorUserID, &entry.Action, &entry.TargetType, &entry.TargetID, &entry.Reason, &entry.IPAddress, &entry.CreatedAt); err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}

	return entries, rows.Err()
}
