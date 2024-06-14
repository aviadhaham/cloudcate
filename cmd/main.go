package main

import (
	"log"
	"os"

	aws_search "github.com/aviadhaham/cloudcate/internal/aws_search/search"
	"github.com/aviadhaham/cloudcate/internal/web"
)

func main() {
	profiles, err := aws_search.GetProfiles()
	if err != nil {
		log.Fatalf("Failed to get profiles: %v", err)
	}

	if os.Getenv("PORT") == "" {
		log.Fatalf("PORT env var is not set")
	}

	s := web.NewServer(os.Getenv("PORT"), profiles)
	s.Run()
}
