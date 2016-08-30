package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xelaadryth/fantasy-friends/fantasy"
)

//MatchForm input fields required for running a match
type MatchForm struct {
	Region          string `form:"region" binding:"required"`
	BlueTeamTop     string `form:"blueTeamTop" binding:"required"`
	BlueTeamJungle  string `form:"blueTeamJungle" binding:"required"`
	BlueTeamMid     string `form:"blueTeamMid" binding:"required"`
	BlueTeamBottom  string `form:"blueTeamBottom" binding:"required"`
	BlueTeamSupport string `form:"blueTeamSupport" binding:"required"`
	RedTeamTop      string `form:"redTeamTop" binding:"required"`
	RedTeamJungle   string `form:"redTeamJungle" binding:"required"`
	RedTeamMid      string `form:"redTeamMid" binding:"required"`
	RedTeamBottom   string `form:"redTeamBottom" binding:"required"`
	RedTeamSupport  string `form:"redTeamSupport" binding:"required"`
}

func playMatch(c *gin.Context) {
	var matchForm MatchForm
	err := c.Bind(&matchForm)
	if err != nil {
		invalidHandler(c, http.StatusBadRequest, err)
		return
	}

	matchScore, err := fantasy.PlayMatch(
		matchForm.Region,
		matchForm.BlueTeamTop,
		matchForm.BlueTeamJungle,
		matchForm.BlueTeamMid,
		matchForm.BlueTeamBottom,
		matchForm.BlueTeamSupport,
		matchForm.RedTeamTop,
		matchForm.RedTeamJungle,
		matchForm.RedTeamMid,
		matchForm.RedTeamBottom,
		matchForm.RedTeamSupport,
	)
	if err != nil {
		invalidHandler(c, http.StatusBadRequest, err)
		return
	}

	c.HTML(http.StatusOK, "match.tmpl", gin.H{
		"matchScore": matchScore,
	})
}
