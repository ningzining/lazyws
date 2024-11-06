package server

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

type Server struct {
	engine *gin.Engine
	port   int
}

func New(port int) *Server {
	engine := gin.Default()
	return &Server{
		engine: engine,
		port:   port,
	}
}

func (s *Server) Start() error {
	if err := s.engine.Run(fmt.Sprintf(":%d", s.port)); err != nil {
		return err
	}

	log.Println("Server started at: ", s.port)
	return nil
}

func (s *Server) AddHandler(route string, handler func(c *gin.Context)) {
	s.engine.GET(route, handler)
}
