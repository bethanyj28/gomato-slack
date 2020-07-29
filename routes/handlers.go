package routes

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nlopes/slack"

	"github.com/bethanyj28/gomato-slack/internal/format"
)

type timerDetail struct {
	timeStart    time.Time
	timeDuration time.Duration
	timer        *time.Timer
}

func (s *Server) startTimer(c *gin.Context) {
	sc, err := slack.SlashCommandParse(c.Request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error parsing slash command": err.Error()})
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
		return
	}

	if _, err := s.Gomato.Start(sc.UserID, d, s.setTimer(sc.UserID)); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error starting timer": err.Error()})
		return
	}

	interpolate := struct {
		TimeDuration string
	}{sc.Text}
	resp, err := format.Message("start", interpolate)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error generating message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (s *Server) pauseTimer(c *gin.Context) {
	sc, err := slack.SlashCommandParse(c.Request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error parsing slash command": err.Error()})
		return
	}

	if !sc.ValidateToken(os.Getenv("SLACK_VERIFICATION_TOKEN")) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if err := s.Gomato.Pause(sc.UserID); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error pausing timer": err.Error()})
		return
	}

	resp, err := format.Message("pause", nil)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error generating message": err.Error()})
	}

	c.JSON(http.StatusOK, resp)
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

	resp, err := format.Message("resume", nil)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error generating message": err.Error()})
	}

	c.JSON(http.StatusOK, resp)
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

	resp, err := format.Message("stop", nil)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error generating message": err.Error()})
	}

	c.JSON(http.StatusOK, resp)
}

func (s *Server) notifyUser(userID string) {
	fmt.Printf("Time's up!")
}

func (s *Server) setTimer(userID string) func() {
	return func() {
		s.notifyUser(userID)
	}
}
