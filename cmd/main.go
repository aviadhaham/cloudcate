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

	if os.Getenv("PORT") == "" {
		log.Fatalf("PORT env var is not set")
	}

	s := web.NewServer(os.Getenv("PORT"), profiles, regions)
	s.Run()
}
