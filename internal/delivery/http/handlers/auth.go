package handlers

import (
	"cat-test/internal/delivery/http/middleware"
	"cat-test/internal/domain"
	"cat-test/internal/errors"
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"net/http"
)

type authHandler struct {
	usecase domain.AuthUsecase
}

type AuthLoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (inp AuthLoginInput) Validate() error {
	return validation.ValidateStruct(&inp,
		validation.Field(&inp.Email, validation.Required),
		validation.Field(&inp.Password, validation.Required),
	)
}

func (h authHandler) Login(c *gin.Context) {
	var input AuthLoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := input.Validate(); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": err})
		return
	}

	token, err := h.usecase.Login(c, input.Email, input.Password)
	if err != nil {
		switch err {
		case errors.AuthInvalidEmail, errors.AuthInvalidPassword:
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, validation.Errors{
				"email":    errors.AuthInvalidEmailOrPassword,
				"password": errors.AuthInvalidEmailOrPassword,
			})
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h authHandler) Logout(c *gin.Context) {
	identity := c.MustGet(middleware.AuthIdentityCtxKey).(domain.User)

	if err := h.usecase.Logout(c, identity.ID); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func SetupAuthHandler(router *gin.RouterGroup, usecase domain.AuthUsecase) {
	handler := authHandler{usecase: usecase}

	router.POST("/login", handler.Login)
	router.POST("/logout", middleware.Authenticate(usecase), handler.Logout)
}
