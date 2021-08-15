package handlers

import (
	"cat-test/internal/domain"
	"database/sql"
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"net/http"
	"strconv"
)

type userHandler struct {
	usecase     domain.UserUsecase
	authUsecase domain.AuthUsecase
	validator   domain.UserValidator
}

func (h userHandler) List(c *gin.Context) {
	users, err := h.usecase.List(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h userHandler) Get(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.usecase.Get(c, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

type UserCreateInput struct {
	Email     string  `json:"email"`
	Password  string  `json:"password"`
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
	Type      int     `json:"type"`
}

func (inp UserCreateInput) Validate(validator domain.UserValidator) error {
	return validation.ValidateStruct(&inp,
		validation.Field(&inp.Email, validation.Required, validation.By(validator.EmailIsUnique)),
		validation.Field(&inp.Password, validation.Required, validation.RuneLength(8, 32)),
		validation.Field(&inp.FirstName, validation.When(inp.FirstName != nil, validation.RuneLength(1, 255))),
		validation.Field(&inp.LastName, validation.When(inp.LastName != nil, validation.RuneLength(1, 255))),
		validation.Field(&inp.Type, validation.Required, validation.In(domain.UserType, domain.UserTypeAdmin)),
	)
}

func (inp UserCreateInput) user() domain.User {
	return domain.User{
		Email:     inp.Email,
		FirstName: inp.FirstName,
		LastName:  inp.LastName,
		Type:      inp.Type,
	}
}

func (h userHandler) Create(c *gin.Context) {
	var input UserCreateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := input.Validate(h.validator); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, err)
		return
	}

	passwordHash, err := h.authUsecase.HashPassword(input.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user := input.user()
	user.PasswordHash = passwordHash

	userID, err := h.usecase.Create(c, user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": userID})
}

type UserUpdateInput struct {
	Email     string  `json:"email"`
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
}

func (inp UserUpdateInput) ValidateAndLoad(validator domain.UserValidator, user *domain.User) error {
	errors := make(validation.Errors)

	if inp.Email != "" && inp.Email != user.Email {
		errors["email"] = validation.Validate(inp.Email, validation.By(validator.EmailIsUnique))
		user.Email = inp.Email
	}

	if inp.FirstName != nil {
		errors["firstName"] = validation.Validate(validation.RuneLength(1, 255))
		user.FirstName = inp.FirstName
	}

	if inp.LastName != nil {
		errors["lastName"] = validation.Validate(validation.RuneLength(1, 255))
		user.LastName = inp.LastName
	}

	return errors.Filter()
}

func (h userHandler) Update(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.usecase.Get(c, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var input UserUpdateInput
	if err = c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err = input.ValidateAndLoad(h.validator, &user); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, err)
		return
	}

	if err = h.usecase.Update(c, user); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h userHandler) Delete(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err = h.usecase.Delete(c, userID); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func SetupUserHandler(router *gin.RouterGroup, usecase domain.UserUsecase, authUsecase domain.AuthUsecase, validator domain.UserValidator) {
	handler := userHandler{usecase: usecase, authUsecase: authUsecase, validator: validator}

	router.GET("", handler.List)
	router.POST("", handler.Create)
	router.GET(":id", handler.Get)
	router.PATCH(":id", handler.Update)
	router.DELETE(":id", handler.Delete)
}
