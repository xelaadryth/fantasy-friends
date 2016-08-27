package controller

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/TrevorSStone/goriot"
	"github.com/gin-gonic/gin"
	"github.com/xelaadryth/fantasy-friends/fantasy"
)

//MatchForm input fields required for running a match
type MatchForm struct {
	Region       string `form:"region" binding:"required"`
	OrderTop     string `form:"orderTop" binding:"required"`
	OrderJungle  string `form:"orderJungle" binding:"required"`
	OrderMid     string `form:"orderMid" binding:"required"`
	OrderBottom  string `form:"orderBottom" binding:"required"`
	OrderSupport string `form:"orderSupport" binding:"required"`
	ChaosTop     string `form:"chaosTop" binding:"required"`
	ChaosJungle  string `form:"chaosJungle" binding:"required"`
	ChaosMid     string `form:"chaosMid" binding:"required"`
	ChaosBottom  string `form:"chaosBottom" binding:"required"`
	ChaosSupport string `form:"chaosSupport" binding:"required"`
}

//TODO: Move most of this code into fantasy
func playMatch(c *gin.Context) {
	var matchForm MatchForm
	err := c.Bind(&matchForm)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
			"form":   matchForm,
		})
		return
	}
	summonerNames := goriot.NormalizeSummonerName(
		matchForm.OrderTop,
		matchForm.OrderJungle,
		matchForm.OrderMid,
		matchForm.OrderBottom,
		matchForm.OrderSupport,
		matchForm.ChaosTop,
		matchForm.ChaosJungle,
		matchForm.ChaosMid,
		matchForm.ChaosBottom,
		matchForm.ChaosSupport,
	)
	if len(summonerNames) != fantasy.PlayersPerMatch {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  fmt.Sprint("Provided ", len(summonerNames), " summoners instead of ", fantasy.PlayersPerMatch, "."),
		})
		return
	}

	//Attempt to get the summoner IDs
	summoners, err := goriot.SummonerByName(matchForm.Region, summonerNames...)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}
	//Make sure the names are unique
	if len(summoners) != fantasy.PlayersPerMatch {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error": fmt.Sprint("Only ", len(summoners), " valid distinct summoner names instead of ",
				fantasy.PlayersPerMatch, "."),
			"summoners": summoners,
		})
		return
	}

	summonerIDs := make([]int64, fantasy.PlayersPerMatch, fantasy.PlayersPerMatch)
	for i := 0; i < fantasy.PlayersPerMatch; i++ {
		summonerIDs[i] = summoners[summonerNames[i]].ID
	}
	log.Println("Calculating scores for ", summonerIDs)
	matchScore, err := fantasy.CalculateScores(matchForm.Region, summonerIDs)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"matchScore": matchScore,
	})
}

//Route does all the routing for the app
func Route() {
	//Get port number to listen for
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	//Set up routing
	router := gin.Default()

	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	router.POST("/matchResults", playMatch)

	router.Run(":" + port)
}
