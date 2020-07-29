package routes

import (
	"github.com/bethanyj28/gomato"
	"github.com/gin-gonic/gin"
)

// Server represents the necessary components to run the server
type Server struct {
	Gomato gomato.PomodoroManager
	Router *gin.Engine
}

// NewServer instantiates a Server object
func NewServer() *Server {
	gm := gomato.NewDefaultTimeKeeper()
	router := gin.Default()
	return &Server{Gomato: gm, Router: router}
}

// BuildRoutes builds the routes to listen to slash commands
func (s *Server) BuildRoutes() {
	timer := s.Router.Group("/timer")
	{
		timer.POST("start", s.startTimer)
		timer.POST("pause", s.pauseTimer)
		timer.POST("resume", s.resumeTimer)
		timer.POST("stop", s.stopTimer)
	}
}
