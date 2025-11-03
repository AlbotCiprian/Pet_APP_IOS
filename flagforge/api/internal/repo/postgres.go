package repo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"

	"github.com/flagforge/flagforge/api/internal/models"
)

// Store encapsulates DB + cache access.
type Store struct {
	db    *pgxpool.Pool
	cache *redis.Client
	log   zerolog.Logger
}

// NewStore constructs a Store.
func NewStore(db *pgxpool.Pool, cache *redis.Client, log zerolog.Logger) *Store {
	return &Store{db: db, cache: cache, log: log}
}

const flagCacheKey = "flags:%s:%s"

func encodeFlagValues(values []models.FlagValue) ([]byte, error) {
	return json.Marshal(values)
}

func decodeFlagValues(payload []byte) ([]models.FlagValue, error) {
	var values []models.FlagValue
	if err := json.Unmarshal(payload, &values); err != nil {
		return nil, err
	}
	return values, nil
}

// FetchFlags returns flag payload for a project/environment pair.
func (s *Store) FetchFlags(ctx context.Context, projectID, env string) ([]models.FlagValue, string, error) {
	key := fmt.Sprintf(flagCacheKey, projectID, env)
	etag, etagErr := s.cache.Get(ctx, key+":etag").Result()
	cached, cacheErr := s.cache.Get(ctx, key).Bytes()
	if etagErr == nil && cacheErr == nil {
		values, decErr := decodeFlagValues(cached)
		if decErr == nil {
			return values, etag, nil
		}
		s.log.Warn().Err(decErr).Msg("failed to decode cache entry")
	}

	rows, err := s.db.Query(ctx, `
SELECT f.id, fv.environment_id, f.key, f.type, fv.value_json, fv.rollout, fv.rules_json, fv.created_at, fv.created_by
FROM flags f
JOIN flag_versions fv ON fv.flag_id = f.id
JOIN environments e ON e.id = fv.environment_id
WHERE f.project_id = $1 AND e.key = $2 AND fv.created_at = (
    SELECT MAX(created_at) FROM flag_versions fv2 WHERE fv2.flag_id = f.id AND fv2.environment_id = fv.environment_id
)
ORDER BY f.key
`, projectID, env)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	values := make([]models.FlagValue, 0)
	for rows.Next() {
		var fv models.FlagValue
		if err := rows.Scan(&fv.FlagID, &fv.EnvironmentID, &fv.Key, &fv.Type, &fv.ValueJSON, &fv.Rollout, &fv.RulesJSON, &fv.UpdatedAt, &fv.UpdatedBy); err != nil {
			return nil, "", err
		}
		values = append(values, fv)
	}
	if rows.Err() != nil {
		return nil, "", rows.Err()
	}

	payload, err := encodeFlagValues(values)
	if err == nil {
		etag = fmt.Sprintf("W/\"%d\"", time.Now().Unix())
		if err := s.cache.Set(ctx, key, payload, 5*time.Minute).Err(); err != nil {
			s.log.Warn().Err(err).Msg("failed to set cache entry")
		}
		if err := s.cache.Set(ctx, key+":etag", etag, 5*time.Minute).Err(); err != nil {
			s.log.Warn().Err(err).Msg("failed to set etag cache")
		}
	}

	return values, etag, nil
}

// CreateFlag inserts a flag definition and seed version.
func (s *Store) CreateFlag(ctx context.Context, f models.Flag, envValues map[string][]byte) error {
	if len(envValues) == 0 {
		return errors.New("envValues cannot be empty")
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
INSERT INTO flags (id, project_id, key, type, description, created_by, created_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
`, f.ID, f.ProjectID, f.Key, f.Type, f.Description, f.CreatedBy, f.CreatedAt)
	if err != nil {
		return err
	}

	for envID, val := range envValues {
		_, err = tx.Exec(ctx, `
INSERT INTO flag_versions (id, flag_id, environment_id, value_json, rollout, rules_json, created_at, created_by)
VALUES (gen_random_uuid(), $1, $2, $3, 100, '{}', NOW(), $4)
`, f.ID, envID, val, f.CreatedBy)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	s.invalidateFlag(ctx, f.ProjectID)
	return nil
}

// ListAuditLogs returns logs for a flag.
func (s *Store) ListAuditLogs(ctx context.Context, flagID string) ([]models.AuditLog, error) {
	rows, err := s.db.Query(ctx, `
SELECT id, actor_id, org_id, entity_type, entity_id, action, diff_json, ts
FROM audit_logs WHERE entity_id = $1
ORDER BY ts DESC
LIMIT 100
`, flagID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	logs := make([]models.AuditLog, 0)
	for rows.Next() {
		var log models.AuditLog
		if err := rows.Scan(&log.ID, &log.ActorID, &log.OrgID, &log.Entity, &log.EntityID, &log.Action, &log.DiffJSON, &log.Timestamp); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	return logs, rows.Err()
}

func (s *Store) invalidateFlag(ctx context.Context, projectID string) {
	pattern := fmt.Sprintf(flagCacheKey, projectID, "*")
	iter := s.cache.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		if err := s.cache.Del(ctx, iter.Val()).Err(); err != nil {
			s.log.Warn().Err(err).Msg("failed to delete cache key")
		}
	}
	if err := iter.Err(); err != nil {
		s.log.Warn().Err(err).Msg("scan error")
	}
	if err := s.cache.Publish(ctx, "flagforge.invalidate", projectID).Err(); err != nil {
		s.log.Warn().Err(err).Msg("publish invalidate")
	}
}
