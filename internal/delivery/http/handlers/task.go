package handlers

import (
	"cat-test/internal/domain"
	"cat-test/internal/errors"
	"database/sql"
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"net/http"
	"strconv"
)

type taskHandler struct {
	usecase   domain.TaskUsecase
	validator domain.TaskValidator
}

func (h taskHandler) List(c *gin.Context) {
	tasks, err := h.usecase.List(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}
func (h taskHandler) Get(c *gin.Context) {
	taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := h.usecase.Get(c, taskID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "task not found"})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
}

type TaskCreateInput struct {
	Title         string `json:"title"`
	Description   string `json:"description"`
	EstimatedTime *int64 `json:"estimatedTime"`
	UserID        *int64 `json:"userId"`
	DueDate       *int64 `json:"dueDate"`
}

func (inp TaskCreateInput) Validate(validator domain.TaskValidator) error {
	return validation.ValidateStruct(&inp,
		validation.Field(&inp.Title, validation.Required, validation.RuneLength(1, 255)),
		validation.Field(&inp.Description, validation.When(inp.Description != "", validation.RuneLength(1, 4096))),
		validation.Field(&inp.UserID, validation.When(inp.UserID != nil, validation.By(validator.ValidateUserExist))),
	)
}

func (inp TaskCreateInput) task() domain.Task {
	return domain.Task{
		Title:         inp.Title,
		Description:   inp.Description,
		EstimatedTime: inp.EstimatedTime,
		UserID:        inp.UserID,
	}
}

func (h taskHandler) Create(c *gin.Context) {
	var input TaskCreateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := input.Validate(h.validator); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, err)
		return
	}

	taskID, err := h.usecase.Create(c, input.task())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": taskID})
}

type TaskUpdateInput struct {
	Title         string `json:"title"`
	Description   string `json:"description"`
	EstimatedTime *int64 `json:"estimatedTime"`
	Status        int    `json:"status"`
	DueDate       *int64 `json:"dueDate"`
}

func (inp TaskUpdateInput) ValidateAndLoad(task *domain.Task) error {
	errs := make(validation.Errors)

	if inp.Title != "" && inp.Title != task.Title {
		errs["title"] = validation.Validate(inp.Title, validation.RuneLength(1, 255))
		task.Title = inp.Title
	}

	if inp.Description != "" && inp.Description != task.Description {
		errs["description"] = validation.Validate(inp.Description, validation.RuneLength(0, 4096))
		task.Description = inp.Description
	}

	if inp.EstimatedTime != nil {
		task.EstimatedTime = inp.EstimatedTime
	}

	if inp.Status != 0 {
		errs["status"] = validation.Validate(inp.Status,
			validation.In(domain.TaskStatusWait, domain.TaskStatusInProcess, domain.TaskStatusCompleted))
		task.Status = inp.Status
	}

	if inp.DueDate != nil {
		task.DueDate = inp.DueDate
	}

	return errs.Filter()
}

func (h taskHandler) Update(c *gin.Context) {
	taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := h.usecase.Get(c, taskID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "task not found"})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var input TaskUpdateInput
	if err = c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err = input.ValidateAndLoad(&task); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, err)
		return
	}

	if err = h.usecase.Update(c, task); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h taskHandler) Delete(c *gin.Context) {
	taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err = h.usecase.Delete(c, taskID); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

type TaskAssignUserInput struct {
	UserID int64 `json:"userId"`
}

func (inp TaskAssignUserInput) Validate(validator domain.TaskValidator) error {
	return validation.ValidateStruct(&inp,
		validation.Field(&inp.UserID, validation.Required, validation.By(validator.ValidateUserExist)),
	)
}

func (h taskHandler) AssignUser(c *gin.Context) {
	taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var input TaskAssignUserInput
	if err = c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err = input.Validate(h.validator); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": err})
		return
	}

	if err = h.usecase.AssignUser(c, taskID, input.UserID); err != nil {
		if err == errors.TaskCannotBeAssigned {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "task has been successfully assigned to the user"})
}

func SetupTaskHandler(router *gin.RouterGroup, usecase domain.TaskUsecase, authUsecase domain.AuthUsecase, validator domain.TaskValidator) {
	handler := taskHandler{usecase: usecase, validator: validator}

	router.GET("", handler.List)
	router.POST("", handler.Create)
	router.GET(":id", handler.Get)
	router.PATCH(":id", handler.Update)
	router.DELETE(":id", handler.Delete)
	router.POST(":id/assign", handler.AssignUser)
}
