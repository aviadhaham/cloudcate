package web

import (
	"log"
	"net/http"

	config "github.com/aviadhaham/cloudcate-service/internal/aws_search/config"
	search "github.com/aviadhaham/cloudcate-service/internal/aws_search/search"

	"github.com/gin-gonic/gin"
)

func NewRouter(profiles []string, regions []string) *gin.Engine {
	r := gin.Default()

	r.Use(CORS())

	// Serve static files from the "static" directory
	r.Static("/static", "./static")

	r.GET("/", func(c *gin.Context) {
		c.File("static/index.html")
	})

	r.GET("/search", func(c *gin.Context) {
		resourceName := c.Query("resource_name")
		resourceType := c.Query("resource_type")

		results, err := search.FindResources(profiles, regions, config.ResourceGlobality, resourceType, resourceName)
		if err != nil {
			log.Fatalf("Failed to search resources: %v", err)
		}

		c.JSON(http.StatusOK, gin.H{
			"results": results,
		})
	})

	return r
}
