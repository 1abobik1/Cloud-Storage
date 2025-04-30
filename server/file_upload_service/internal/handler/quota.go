package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// InitUserPlan создаёт для userID план free
func (h *Handler) InitUserPlan(c *gin.Context) {
	idStr := c.Param("id")
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	if err := h.quotaService.InitializeFreePlan(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}
