package routes

import (
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func loadDocumentationRoutes(group *echo.Group) {
	group.GET("/documentation/*", echoSwagger.WrapHandler)
}
