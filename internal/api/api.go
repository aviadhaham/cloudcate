package api

import (
	"log"
	"net/http"

	"github.com/aviadhaham/cloudcate/internal/config"
	"github.com/aviadhaham/cloudcate/internal/search"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type server struct {
	router *echo.Echo
}

func NewServer(profiles []string) *server {
	s := &server{}
	s.router = echo.New()
	s.router.Use(middleware.Logger())
	s.router.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:   "./web/dist",
		Browse: true,
	}))
	apiGroup := s.router.Group("/api")
	apiGroup.GET("/search", func(c echo.Context) error {
		resourceName := c.QueryParam("resource_name")
		resourceType := c.QueryParam("resource_type")
		resourceSubType := c.QueryParam("resource_subtype")

		results, err := search.FindResources(profiles, config.ServicesGlobality, resourceType, resourceSubType, resourceName)
		if err != nil {
			log.Fatalf("Failed to search resources: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to search resources"})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"results": results,
		})
	})

	return s
}

func (s *server) Run(port string) {
	s.router.Logger.Fatal(s.router.Start(":" + port))
}
