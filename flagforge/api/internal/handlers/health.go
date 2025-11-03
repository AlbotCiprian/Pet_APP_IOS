package handlers

import "github.com/gin-gonic/gin"

// Health responds with service availability.
func (h Handler) Health(c *gin.Context) {
	c.JSON(200, gin.H{"ok": true})
}
