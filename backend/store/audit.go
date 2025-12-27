package store

import (
	"encoding/json"
	"time"

	uuid "github.com/satori/go.uuid"
)

func LogAudit(db *DB, actorID, actorEmail, action, entityType, entityID, tenantID string, details map[string]interface{}) error {
	if db == nil || db.Conn == nil {
		return nil
	}
	id := uuid.NewV4().String()
	var payload []byte
	if details != nil {
		if b, err := json.Marshal(details); err == nil {
			payload = b
		}
	}
	_, err := db.Conn.Exec(
		`INSERT INTO audit_logs (id, actor_user_id, actor_email, tenant_id, action, entity_type, entity_id, details, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		id, actorID, actorEmail, nullableString(tenantID), action, entityType, nullableString(entityID), payload, time.Now(),
	)
	return err
}

func nullableString(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}
