package gomato

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nlopes/slack"
)

type timerDetail struct {
	timeStart    time.Time
	timeDuration time.Duration
	timer        *time.Timer
}

func (s *Server) startTimer(c *gin.Context) {
	sc, err := slack.SlashCommandParse(c.Request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error parsing slash command": err.Error()})
		return
	}

	if !sc.ValidateToken(os.Getenv("SLACK_VERIFICATION_TOKEN")) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	ud := sc.Text
	if ud == "" {
		ud = "20"
	}

	ud = fmt.Sprintf("%sm", ud)

	d, err := time.ParseDuration(ud)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	if _, err := s.Gomato.Start(sc.UserID, d, s.setTimer(sc.UserID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error starting timer": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (s *Server) pauseTimer(c *gin.Context) {
	sc, err := slack.SlashCommandParse(c.Request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error parsing slash command": err.Error()})
		return
	}

	if !sc.ValidateToken(os.Getenv("SLACK_VERIFICATION_TOKEN")) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if err := s.Gomato.Pause(sc.UserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error pausing timer": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (s *Server) resumeTimer(c *gin.Context) {
	sc, err := slack.SlashCommandParse(c.Request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error parsing slash command": err.Error()})
		return
	}

	if !sc.ValidateToken(os.Getenv("SLACK_VERIFICATION_TOKEN")) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if err := s.Gomato.Resume(sc.UserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error resuming timer": err.Error()})
	}

	c.Status(http.StatusOK)
}

func (s *Server) stopTimer(c *gin.Context) {
	sc, err := slack.SlashCommandParse(c.Request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error parsing slash command": err.Error()})
		return
	}

	if !sc.ValidateToken(os.Getenv("SLACK_VERIFICATION_TOKEN")) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if err := s.Gomato.Stop(sc.UserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error stopping timer": err.Error()})
	}

	c.Status(http.StatusOK)
}

func (s *Server) notifyUser(userID string) {
	fmt.Printf("Time's up!")
}

func (s *Server) setTimer(userID string) func() {
	return func() {
		s.notifyUser(userID)
	}
}
