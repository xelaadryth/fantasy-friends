package controller

import (
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/xelaadryth/fantasy-friends/fantasy"
)

func playMatch(c *gin.Context) {
	var blueTeamForm BlueTeamForm
	var redTeamForm RedTeamForm
	err := c.Bind(&blueTeamForm)
	if err != nil {
		invalidHandler(c, http.StatusBadRequest, err)
		return
	}
	err = c.Bind(&redTeamForm)
	if err != nil {
		invalidHandler(c, http.StatusBadRequest, err)
		return
	}

	matchScore, err := fantasy.PlayMatch(
		blueTeamForm.Region,
		blueTeamForm.BlueTeamTop,
		blueTeamForm.BlueTeamJungle,
		blueTeamForm.BlueTeamMid,
		blueTeamForm.BlueTeamBottom,
		blueTeamForm.BlueTeamSupport,
		redTeamForm.RedTeamTop,
		redTeamForm.RedTeamJungle,
		redTeamForm.RedTeamMid,
		redTeamForm.RedTeamBottom,
		redTeamForm.RedTeamSupport,
	)
	if err != nil {
		invalidHandler(c, http.StatusBadRequest, err)
		return
	}
	session := sessions.Default(c)
	session.Set(sessionNavActive, "play")
	sessionMap := sessionAsMap(&session)

	c.HTML(http.StatusOK, "match.tmpl", gin.H{
		sessionName:  *sessionMap,
		"matchScore": matchScore,
	})
}
