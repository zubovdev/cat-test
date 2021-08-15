package http

import (
	"cat-test/internal/delivery/http/handlers"
	"cat-test/internal/delivery/http/middleware"
	"cat-test/internal/usecase"
	"cat-test/internal/validator"
	"github.com/gin-gonic/gin"
)

func getRouter(usecases usecase.Usecases, validators validator.Validators) *gin.Engine {
	router := gin.Default()

	handlers.SetupAuthHandler(router.Group("/api/auth"), usecases.Auth)

	apiRouter := router.Group("/api")
	apiRouter.Use(middleware.Authenticate(usecases.Auth))
	handlers.SetupUserHandler(apiRouter.Group("/users"), usecases.User, usecases.Auth, validators.User)
	handlers.SetupTaskHandler(apiRouter.Group("/tasks"), usecases.Task, usecases.Auth, validators.Task)

	return router
}
