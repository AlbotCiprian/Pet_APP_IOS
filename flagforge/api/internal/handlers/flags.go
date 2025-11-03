package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"github.com/flagforge/flagforge/api/internal/models"
	"github.com/flagforge/flagforge/api/internal/repo"
)

// Handler bundles dependencies for HTTP handlers.
type Handler struct {
	Store *repo.Store
	Log   zerolog.Logger
}

// ListFlags returns the current flag values for a project and environment.
func (h Handler) ListFlags(c *gin.Context) {
	projectID := c.Query("project_id")
	env := c.Query("env")
	if projectID == "" || env == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "project_id and env are required"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	values, etag, err := h.Store.FetchFlags(ctx, projectID, env)
	if err != nil {
		h.Log.Error().Err(err).Msg("fetch flags")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to fetch flags"})
		return
	}

	if match := c.GetHeader("If-None-Match"); match != "" && etag != "" && match == etag {
		c.Status(http.StatusNotModified)
		return
	}

	if etag != "" {
		c.Header("ETag", etag)
	}
	if values == nil {
		values = []models.FlagValue{}
	}
	c.JSON(http.StatusOK, gin.H{"flags": values})
}

// CreateFlag registers a new feature flag placeholder.
func (h Handler) CreateFlag(c *gin.Context) {
	var req struct {
		ProjectID string            `json:"project_id" binding:"required"`
		Key       string            `json:"key" binding:"required"`
		Type      string            `json:"type" binding:"required"`
		Values    map[string]string `json:"values" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	flag := models.Flag{
		ID:          uuid.NewString(),
		ProjectID:   req.ProjectID,
		Key:         req.Key,
		Type:        req.Type,
		Description: "",
		CreatedBy:   "system",
		CreatedAt:   time.Now().UTC(),
	}

	envValues := make(map[string][]byte, len(req.Values))
	for envID, val := range req.Values {
		envValues[envID] = []byte(val)
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	if err := h.Store.CreateFlag(ctx, flag, envValues); err != nil {
		h.Log.Error().Err(err).Msg("create flag")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to create flag"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": flag.ID})
}

// FlagAudit returns audit events for a flag.
func (h Handler) FlagAudit(c *gin.Context) {
	flagID := c.Param("id")
	if flagID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing id"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	logs, err := h.Store.ListAuditLogs(ctx, flagID)
	if err != nil {
		h.Log.Error().Err(err).Msg("list audit logs")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to fetch audit logs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"logs": logs})
}
