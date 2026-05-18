package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"tasks-manager/internal/domain"
	"tasks-manager/internal/service"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	service service.TaskService
}

func NewHandler(service service.TaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

func (h *TaskHandler) Create(c *gin.Context) {
	var task domain.Task

	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _, err := ValidateUserContext(c)
	if err != nil {
		if err.Error() == "invalid or expired token" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}
		if err.Error() == "internal server error" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}
	}

	if err := h.service.Create(&task, int64(userID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, task)
}

func (h *TaskHandler) GetById(c *gin.Context) {
	userID, isAdmin, err := ValidateUserContext(c)
	if err != nil {
		if err.Error() == "invalid or expired token" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}
		if err.Error() == "internal server error" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}
	}

	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	task, err := h.service.GetById(id, userID, isAdmin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) GetAll(c *gin.Context) {
	userID, _, err := ValidateUserContext(c)
	if err != nil {
		if err.Error() == "invalid or expired token" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}
		if err.Error() == "internal server error" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}
	}

	tasks, err := h.service.GetAll(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func (h *TaskHandler) Update(c *gin.Context) {
	userID, _, err := ValidateUserContext(c)
	if err != nil {
		if err.Error() == "invalid or expired token" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}
		if err.Error() == "internal server error" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}
	}
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	var task domain.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task.ID = id

	if err := h.service.Update(&task, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) Delete(c *gin.Context) {
	userID, isAdmin, err := ValidateUserContext(c)
	if err != nil {
		if err.Error() == "invalid or expired token" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}
		if err.Error() == "internal server error" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}
	}

	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.Delete(id, userID, isAdmin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)

}

func ValidateUserContext(c *gin.Context) (int64, bool, error) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		return 0, false, fmt.Errorf("invalid or expired token")
	}

	userID, ok := userIDValue.(int64)
	if !ok {
		return 0, false, fmt.Errorf("internal server error")
	}

	isAdmin := false
	userRoleValue, exists := c.Get("role")
	if !exists {
		return 0, false, fmt.Errorf("invalid or expired token")
	}
	userRole, ok := userRoleValue.(string)
	if !ok {
		return 0, false, fmt.Errorf("internal server error")
	}
	if userRole == "admin" {
		isAdmin = true
	}

	return userID, isAdmin, nil
}
