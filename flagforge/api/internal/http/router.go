package http

import (
	"github.com/gin-gonic/gin"

	"github.com/flagforge/flagforge/api/internal/handlers"
)

// NewRouter wires HTTP routes.
func NewRouter(handler handlers.Handler) *gin.Engine {
	r := gin.Default()

	r.GET("/healthz", handler.Health)

	v1 := r.Group("/v1")
	{
		v1.GET("/flags", handler.ListFlags)
		v1.POST("/flags", handler.CreateFlag)
		v1.GET("/flags/:id/audit", handler.FlagAudit)
	}

	return r
}
