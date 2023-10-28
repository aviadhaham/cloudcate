package main

import (
	"log"
	"os"

	aws_search "github.com/aviadhaham/cloudcate-service/internal/aws_search/search"
	"github.com/aviadhaham/cloudcate-service/internal/web"
)

func main() {
	profiles, err := aws_search.GetProfiles()
	if err != nil {
		log.Fatalf("Failed to get profiles: %v", err)
	}

	regions, err := aws_search.GetRegions(profiles[0])
	if err != nil {
		log.Fatalf("Failed to get regions: %v", err)
	}

	if os.Getenv("PORT") == "" {
		log.Fatalf("PORT env var is not set")
	}

	s := web.NewServer(os.Getenv("PORT"), profiles, regions)
	s.Run()
}
