package controller

import (
	"errors"
	"log"
	"net/http"

	"github.com/TrevorSStone/goriot"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/xelaadryth/fantasy-friends/database"
	"github.com/xelaadryth/fantasy-friends/fantasy"
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

	userID, err := validateSession(&session)
	if err != nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	//TODO: Update region to be associated with individual players rather than user account session
	regionUntyped := session.Get(sessionRegion)
	region, ok := regionUntyped.(string)
	if !ok {
		region = "na"
	}

	//Get the team from the DB
	teams, _ := database.GetTeams(userID)
	//TODO: Be able to handle more than one team
	if len(*teams) > 0 {
		team := (*teams)[0]
		//Convert summoner IDs to summoner names
		summonersMap, err := fantasy.GetSummonersByID(
			//TODO: Create our own summoner DB that allows cross-region teams
			region,
			team.Top,
			team.Jungle,
			team.Mid,
			team.Bottom,
			team.Support,
		)
		if err != nil {
			database.DeleteTeam(team.ID)
			invalidHandler(c, http.StatusBadRequest, errors.New("Stored summoner IDs invalid, clearing the team."))
			return
		}

		session.Set(sessionTeam, map[string]string{
			sessionTop:     summonersMap[team.Top].Name,
			sessionJungle:  summonersMap[team.Jungle].Name,
			sessionMid:     summonersMap[team.Mid].Name,
			sessionBottom:  summonersMap[team.Bottom].Name,
			sessionSupport: summonersMap[team.Support].Name,
			sessionName:    team.Name,
		})
	} else {

		session.Set(sessionTeam, map[string]string{
			sessionTop:     "",
			sessionJungle:  "",
			sessionMid:     "",
			sessionBottom:  "",
			sessionSupport: "",
			sessionName:    "New Team",
		})
	}

	session.Set(sessionNavActive, "team")
	sessionMap := sessionAsMap(&session)
	c.HTML(http.StatusOK, "team.tmpl", gin.H{
		sessionSession: *sessionMap,
	})
}

//saveTeam to player in the database
func saveTeam(c *gin.Context) {
	session := sessions.Default(c)

	var blueTeamForm BlueTeamForm
	err := c.Bind(&blueTeamForm)
	if err != nil {
		invalidHandler(c, http.StatusBadRequest, err)
		return
	}

	userID, err := validateSession(&session)
	if err != nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}
	username, err := database.GetUsername(userID)
	if err != nil {
		invalidHandler(c, http.StatusBadRequest, errors.New("User account missing."))
		return
	}

	//TODO: This should be kept on a per-player basis, not a per-user session basis
	session.Set(sessionRegion, blueTeamForm.Region)
	normalizedNames := goriot.NormalizeSummonerName(
		blueTeamForm.BlueTeamTop,
		blueTeamForm.BlueTeamJungle,
		blueTeamForm.BlueTeamMid,
		blueTeamForm.BlueTeamBottom,
		blueTeamForm.BlueTeamSupport,
	)
	summonersMap, err := fantasy.GetSummonersByName(blueTeamForm.Region, normalizedNames...,
	)
	if err != nil {
		invalidHandler(c, http.StatusBadRequest, err)
		return
	}

	teams, err := database.GetTeams(userID)
	if err != nil {
		log.Println("Error retrieving teams for userID", userID)
	}

	//Team created from the form info
	//TODO: Allow more than one team
	team := database.FantasyTeam{
		Owner: userID,
		//TODO: Allow modifications to team name
		Name: username + "'s Team",
		//TODO: Allow display order swapping once we can have more than one team
		Position: 1,
		Top:      summonersMap[normalizedNames[0]].ID,
		Jungle:   summonersMap[normalizedNames[1]].ID,
		Mid:      summonersMap[normalizedNames[2]].ID,
		Bottom:   summonersMap[normalizedNames[3]].ID,
		Support:  summonersMap[normalizedNames[4]].ID,
	}

	if len(*teams) > 0 {
		team.ID = (*teams)[0].ID
		err = database.UpdateTeam(team)
		if err != nil {
			invalidHandler(c, http.StatusBadRequest, errors.New("Error updating team."))
			return
		}
	} else {
		//ID field is ignored
		err = database.AddTeam(team)
		if err != nil {
			invalidHandler(c, http.StatusBadRequest, errors.New("Error creating new team."))
			return
		}
	}

	c.Redirect(http.StatusFound, "/team")
}
