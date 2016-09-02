package controller

import (
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

//BlueTeamForm input fields
type BlueTeamForm struct {
	Region          string `form:"region" binding:"required"`
	BlueTeamTop     string `form:"blueTeamTop" binding:"required"`
	BlueTeamJungle  string `form:"blueTeamJungle" binding:"required"`
	BlueTeamMid     string `form:"blueTeamMid" binding:"required"`
	BlueTeamBottom  string `form:"blueTeamBottom" binding:"required"`
	BlueTeamSupport string `form:"blueTeamSupport" binding:"required"`
}

//RedTeamForm input fields
type RedTeamForm struct {
	Region         string `form:"region" binding:"required"`
	RedTeamTop     string `form:"redTeamTop" binding:"required"`
	RedTeamJungle  string `form:"redTeamJungle" binding:"required"`
	RedTeamMid     string `form:"redTeamMid" binding:"required"`
	RedTeamBottom  string `form:"redTeamBottom" binding:"required"`
	RedTeamSupport string `form:"redTeamSupport" binding:"required"`
}

func routeTeam(c *gin.Context) {
	session := sessions.Default(c)

	err := validateSession(&session)
	if err != nil {
		c.Redirect(http.StatusFound, "/login")
	}

	//TODO: Load data from db

	session.Set(sessionNavActive, "team")
	sessionMap := sessionAsMap(&session)
	c.HTML(http.StatusOK, "team.tmpl", gin.H{
		sessionName: *sessionMap,
		//TODO: Display info
	})
}

func saveTeam(c *gin.Context) {
	var blueTeamForm BlueTeamForm
	err := c.Bind(&blueTeamForm)
	if err != nil {
		invalidHandler(c, http.StatusBadRequest, err)
		return
	}
	//TODO: Save team data to db

	c.Redirect(http.StatusFound, "/team")
}
