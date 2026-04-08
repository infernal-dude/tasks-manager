package handler

import (
	"net/http"
	"strconv"
	"tasks-manager/internal/domain"
	"tasks-manager/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) Register(c *gin.Context) {
	// Ты используешь одну и ту же структуру таска и юзера в разных местах и это неправильно
	// Во-первых, ты используешь их без JSON тегов, это очень грубая ошибка
	// Во-вторых, для HTTP слоя и для слоя БД у тебя должны быть разные структуры,
	// потому в запросе ты не передаёшь например ID и мне нужно очень глубоко вчитываться в код чтобы понять почему,
	// структура в HTTP слое должна явно отражать тело запроса и не иметь ничего лишнего
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Register(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) Login(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error in logging"})
		return
	}

	token, err := h.service.Login(user)
	if err != nil {
		// Почему 400? А если у тебя база упала? Тогда должно быть 500 явно
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Здесь дложно быть 200 успешно
	// При наведении на код всегда есть ссылка на RFC где описывается его значение
	// Вот что написано про Accepted:
	// The 202 (Accepted) status code indicates that the request has been accepted for processing, but the processing has not been completed.
	// ТАк что он точно сюда не подходит
	c.JSON(http.StatusAccepted, gin.H{"token": token})
}

func (h *UserHandler) GetByUsername(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "empty username"})
		return
	}

	user, err := h.service.GetByUsername(username)
	if err != nil {
		// Тоже самое, почему 404? Упала база, должно быть 500
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) GetById(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	user, err := h.service.GetById(id)
	if err != nil {
		// Тоже самое, почему 404? Упала база, должно быть 500
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	if idParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "empty"})
		return
	}

	var user domain.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "err.Error()"})
		return
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	userOld, err := h.service.GetById(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user.ID = userOld.ID

	err = h.service.Update(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	if idParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "empty"})
		return
	}

	var user domain.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	// Зачем? Можно же просто удалить
	if _, err := h.service.GetById(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.service.Delete(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Здесь должно быть 200, успешно
	c.Status(http.StatusNoContent)
}
