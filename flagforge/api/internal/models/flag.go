package models

import (
	"encoding/json"
	"time"
)

// Flag represents a feature flag definition.
type Flag struct {
	ID          string    `json:"id"`
	ProjectID   string    `json:"project_id"`
	Key         string    `json:"key"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
}

// FlagValue represents an evaluated value for a flag in an environment.
type FlagValue struct {
	FlagID        string          `json:"flag_id"`
	EnvironmentID string          `json:"environment_id"`
	Key           string          `json:"key"`
	Type          string          `json:"type"`
	ValueJSON     json.RawMessage `json:"value_json"`
	Rollout       int             `json:"rollout"`
	RulesJSON     json.RawMessage `json:"rules_json"`
	UpdatedAt     time.Time       `json:"updated_at"`
	UpdatedBy     string          `json:"updated_by"`
}

// AuditLog captures changes to entities in the system.
type AuditLog struct {
	ID        string    `json:"id"`
	ActorID   string    `json:"actor_id"`
	OrgID     string    `json:"org_id"`
	Entity    string    `json:"entity_type"`
	EntityID  string    `json:"entity_id"`
	Action    string    `json:"action"`
	DiffJSON  []byte    `json:"diff_json"`
	Timestamp time.Time `json:"ts"`
}
