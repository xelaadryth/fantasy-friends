package controller

import (
	"errors"
	"fmt"
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

//routeHome to home page
func routeHome(c *gin.Context) {
	session := sessions.Default(c)

	var teams *[]database.FantasyTeam
	var myTeams *[]database.FantasyTeam
	//Requires at least 2 teams in the database
	teams, err := database.GetRandomTeams(2)
	if err != nil {
		fmt.Println(err)
		invalidHandler(c, http.StatusBadRequest, errors.New("Failed to pull 2 random teams from database."))
		return
	}
	//Pick random teams from DB, add user's team if logged in
	userID, err := validateSession(&session)
	if err == nil {
		myTeams, err = database.GetTeams(userID)
		if err != nil {
			fmt.Println(err)
			invalidHandler(c, http.StatusBadRequest, errors.New("Failed to retrieve user team from database."))
			return
		}
		if len(*myTeams) > 0 {
			(*teams)[0] = (*myTeams)[0]
		}
	}

	team := (*teams)[0]
	enemyTeam := (*teams)[1]
	//Convert summoner IDs to summoner names
	cacheIDToSummoner, _ := fantasy.GetSummonersByCacheID(
		team.Top,
		team.Jungle,
		team.Mid,
		team.Bottom,
		team.Support,
	)
	enemyCacheIDToSummoner, _ := fantasy.GetSummonersByCacheID(
		enemyTeam.Top,
		enemyTeam.Jungle,
		enemyTeam.Mid,
		enemyTeam.Bottom,
		enemyTeam.Support,
	)

	//TODO: Region per player for this and team page
	session.Set(sessionTeam, map[string]string{
		sessionTop:     cacheIDToSummoner[team.Top].Name,
		sessionJungle:  cacheIDToSummoner[team.Jungle].Name,
		sessionMid:     cacheIDToSummoner[team.Mid].Name,
		sessionBottom:  cacheIDToSummoner[team.Bottom].Name,
		sessionSupport: cacheIDToSummoner[team.Support].Name,
	})
	session.Set(sessionEnemyTeam, map[string]string{
		sessionTop:     enemyCacheIDToSummoner[enemyTeam.Top].Name,
		sessionJungle:  enemyCacheIDToSummoner[enemyTeam.Jungle].Name,
		sessionMid:     enemyCacheIDToSummoner[enemyTeam.Mid].Name,
		sessionBottom:  enemyCacheIDToSummoner[enemyTeam.Bottom].Name,
		sessionSupport: enemyCacheIDToSummoner[enemyTeam.Support].Name,
	})

	session.Set(sessionNavActive, "home")
	sessionMap := sessionAsMap(&session)
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		sessionSession: *sessionMap,
	})
}

func routeTeam(c *gin.Context) {
	session := sessions.Default(c)

	//Can't view this page unless the user is logged in
	userID, err := validateSession(&session)
	if err != nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	//Get the team from the DB
	teams, err := database.GetTeams(userID)
	if err != nil {
		fmt.Println(err)
		invalidHandler(c, http.StatusBadRequest, errors.New("Failed to retrieve team from database."))
		return
	}
	//TODO: Be able to handle more than one team
	if len(*teams) > 0 {
		team := (*teams)[0]
		//Convert summoner IDs to summoner names
		cacheIDToSummoner, err := fantasy.GetSummonersByCacheID(
			team.Top,
			team.Jungle,
			team.Mid,
			team.Bottom,
			team.Support,
		)
		if err != nil {
			fmt.Println(err)
			database.DeleteTeam(team.ID)
			invalidHandler(c, http.StatusBadRequest, errors.New("Stored summoner IDs invalid, clearing the team."))
			return
		}

		session.Set(sessionTeam, map[string]string{
			sessionTop:     cacheIDToSummoner[team.Top].Name,
			sessionJungle:  cacheIDToSummoner[team.Jungle].Name,
			sessionMid:     cacheIDToSummoner[team.Mid].Name,
			sessionBottom:  cacheIDToSummoner[team.Bottom].Name,
			sessionSupport: cacheIDToSummoner[team.Support].Name,
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

	//You need to be logged in to modify your team
	userID, err := validateSession(&session)
	if err != nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}
	username, err := database.GetUsername(userID)
	if err != nil {
		fmt.Println(err)
		invalidHandler(c, http.StatusBadRequest, errors.New("User account missing."))
		return
	}

	//TODO: This should be kept on a per-player basis, not a per-user session basis
	normalizedNames := goriot.NormalizeSummonerName(
		blueTeamForm.BlueTeamTop,
		blueTeamForm.BlueTeamJungle,
		blueTeamForm.BlueTeamMid,
		blueTeamForm.BlueTeamBottom,
		blueTeamForm.BlueTeamSupport,
	)
	_, nameToCacheID, err := fantasy.GetSummonersByName(blueTeamForm.Region, normalizedNames...)
	if err != nil {
		invalidHandler(c, http.StatusBadRequest, err)
		return
	}

	teams, err := database.GetTeams(userID)
	if err != nil {
		log.Println(err)
		log.Println("Error retrieving teams for userID", userID)
	}

	//TODO: Allow more than one team
	if len(*teams) > 0 {
		//TODO: Allow modifications to team name and display position
		err = database.UpdateTeam(
			(*teams)[0].ID,
			userID,
			username+"'s Team",
			1,
			nameToCacheID[normalizedNames[0]],
			nameToCacheID[normalizedNames[1]],
			nameToCacheID[normalizedNames[2]],
			nameToCacheID[normalizedNames[3]],
			nameToCacheID[normalizedNames[4]],
		)
		if err != nil {
			log.Println(err)
			invalidHandler(c, http.StatusBadRequest, errors.New("Error updating team."))
			return
		}
	} else {
		//TODO: Allow modifications to team name and display position
		err = database.AddTeam(
			userID,
			username+"'s Team",
			1,
			nameToCacheID[normalizedNames[0]],
			nameToCacheID[normalizedNames[1]],
			nameToCacheID[normalizedNames[2]],
			nameToCacheID[normalizedNames[3]],
			nameToCacheID[normalizedNames[4]],
		)
		if err != nil {
			log.Println(err)
			invalidHandler(c, http.StatusBadRequest, errors.New("Error creating new team."))
			return
		}
	}

	c.Redirect(http.StatusFound, "/team")
}
