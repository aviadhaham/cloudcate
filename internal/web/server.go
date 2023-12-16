package web

import (
	"log"
)

type server struct {
	port     string
	profiles []string
	regions  []string
}

func NewServer(port string, profiles []string) *server {
	return &server{
		port:     port,
		profiles: profiles,
	}
}

func (s *server) Run() {
	r := NewRouter(s.profiles)

	err := r.Run(":" + s.port)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
