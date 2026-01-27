package acl

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/FLASH2332/Secure-Question-Paper-Distribution-Portal/internal/models"
)

// AccessDeniedError represents an authorization failure
type AccessDeniedError struct {
	Role       string
	ObjectType string
	Action     string
}

func (e *AccessDeniedError) Error() string {
	return fmt.Sprintf("access denied: %s role cannot %s %s", e.Role, e.Action, e.ObjectType)
}

// EnforcePermission checks permission and logs the attempt
func EnforcePermission(db *sql.DB, user *models.User, objectType, action string, objectID *int) error {
	// Check permission
	allowed, err := CheckPermission(db, user.Role, objectType, action)
	if err != nil {
		logAuditEntry(db, user.ID, action, objectType, objectID, false, err.Error())
		return err
	}

	if !allowed {
		err := &AccessDeniedError{
			Role:       user.Role,
			ObjectType: objectType,
			Action:     action,
		}
		logAuditEntry(db, user.ID, action, objectType, objectID, false, err.Error())
		return err
	}

	// Log successful authorization
	logAuditEntry(db, user.ID, action, objectType, objectID, true, "permission granted")
	return nil
}

// logAuditEntry records access attempts to audit log
func logAuditEntry(db *sql.DB, userID int, action, objectType string, objectID *int, success bool, details string) {
	query := `
        INSERT INTO audit_log (user_id, action, object_type, object_id, success, details, ip_address)
        VALUES (?, ?, ?, ?, ?, ?, ?)
    `

	var objID interface{}
	if objectID != nil {
		objID = *objectID
	} else {
		objID = nil
	}

	_, err := db.Exec(query, userID, action, objectType, objID, success, details, "127.0.0.1")
	if err != nil {
		fmt.Printf("⚠️  Warning: Failed to log audit entry: %v\n", err)
	}
}

// GetAuditLog retrieves audit log entries
func GetAuditLog(db *sql.DB, userID int, limit int) ([]AuditEntry, error) {
	query := `
        SELECT id, user_id, action, object_type, object_id, timestamp, success, details
        FROM audit_log
        WHERE user_id = ?
        ORDER BY timestamp DESC
        LIMIT ?
    `

	rows, err := db.Query(query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get audit log: %w", err)
	}
	defer rows.Close()

	var entries []AuditEntry
	for rows.Next() {
		var entry AuditEntry
		var objectID sql.NullInt64

		err := rows.Scan(
			&entry.ID,
			&entry.UserID,
			&entry.Action,
			&entry.ObjectType,
			&objectID,
			&entry.Timestamp,
			&entry.Success,
			&entry.Details,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan audit entry: %w", err)
		}

		if objectID.Valid {
			id := int(objectID.Int64)
			entry.ObjectID = &id
		}

		entries = append(entries, entry)
	}

	return entries, nil
}

// AuditEntry represents a logged access attempt
type AuditEntry struct {
	ID         int
	UserID     int
	Action     string
	ObjectType string
	ObjectID   *int
	Timestamp  time.Time
	Success    bool
	Details    string
}
