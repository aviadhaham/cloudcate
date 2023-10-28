package web

import (
	"log"
)

type server struct {
	port     string
	profiles []string
	regions  []string
}

func NewServer(port string, profiles []string, regions []string) *server {
	return &server{
		port:     port,
		profiles: profiles,
		regions:  regions,
	}
}

func (s *server) Run() {
	r := NewRouter(s.profiles, s.regions)

	err := r.Run(":" + s.port)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
