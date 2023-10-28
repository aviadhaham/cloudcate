package main

import (
	"fmt"
	"log"
	"os"

	general "github.com/aviadhaham/cloudcate-service/internal/aws_search/search/general"
	"github.com/aviadhaham/cloudcate-service/internal/web"
)

func main() {

	profiles, err := general.GetProfiles()
	if err != nil {
		log.Fatalf("Failed to get profiles: %v", err)
	}

	regions, err := general.GetRegions(profiles[0])
	if err != nil {
		log.Fatalf("Failed to get regions: %v", err)
	}

	r := web.NewRouter(profiles, regions)

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatalf("PORT environment variable not set")
	}

	r.Run(fmt.Sprintf(":%s", port))
}
