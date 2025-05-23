package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/1abobik1/Cloud-Storage/file_upload_service/internal/dto"
	"github.com/1abobik1/Cloud-Storage/file_upload_service/internal/services/quota"
	"github.com/gin-gonic/gin"
)

const (
	bytesInGB = 1024 * 1024 * 1024
	bytesInMB = 1024 * 1024
	bytesInKB = 1024
)

// InitUserPlan создаёт для userID план free
func (h *Handler) InitUserPlan(c *gin.Context) {
	idStr := c.Param("id")
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	if err := h.quotaService.InitializeFreePlan(c, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

// GetUserUsage выдает текущее кол-во используемой памяти по id пользователя
func (h *Handler) GetUserUsage(c *gin.Context) {
	idStr := c.Param("id")
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	domainUserUsage, err := h.quotaService.GetUserUsage(c, userID)
	if err != nil {
		if errors.Is(err, quota.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": quota.ErrUserNotFound.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	used := domainUserUsage.CurrentUsed

    // целые гигабайты
    gbCount := used / bytesInGB
    remAfterGB := used % bytesInGB

    // целые мегабайты из остатка после гигабайтов
    mbCount := remAfterGB / bytesInMB
    remAfterMB := remAfterGB % bytesInMB

    // целые килобайты из остатка после мегабайт
    kbCount := remAfterMB / bytesInKB

    resp := dto.UserUsage{
        CurrentUsedGB:  int(gbCount),
        CurrentUsedMB:  int(mbCount),
        CurrentUsedKB:  int(kbCount),
        StorageLimitGB: int(domainUserUsage.StorageLimit / bytesInGB),
        PlanName:       domainUserUsage.PlanName,
    }

	c.JSON(http.StatusOK, resp)
}
