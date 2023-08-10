package routes

import "github.com/labstack/echo/v4"

func LoadRoutes() *echo.Echo {
	router := echo.New()

	apiGroup := router.Group("/api/v1")
	loadDocumentationRoutes(apiGroup)

	return router
}
