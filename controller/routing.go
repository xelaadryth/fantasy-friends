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

func playMatch(c *gin.Context) {
	region := c.PostForm("region")
	summonerNames := goriot.NormalizeSummonerName(
		c.PostForm("orderTop"),
		c.PostForm("orderJungle"),
		c.PostForm("orderMid"),
		c.PostForm("orderBottom"),
		c.PostForm("orderSupport"),
		c.PostForm("chaosTop"),
		c.PostForm("chaosJungle"),
		c.PostForm("chaosMid"),
		c.PostForm("chaosBottom"),
		c.PostForm("chaosSupport"))
	if len(summonerNames) != fantasy.PlayersPerMatch {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  fmt.Sprint("Provided ", len(summonerNames), " summoners instead of ", fantasy.PlayersPerMatch, "."),
		})
	}

	//Attempt to get the summoner IDs
	summoners, err := goriot.SummonerByName(region, summonerNames...)
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
	}

	summonerIDs := make([]int64, fantasy.PlayersPerMatch, fantasy.PlayersPerMatch)
	for i := 0; i < fantasy.PlayersPerMatch; i++ {
		summonerIDs[i] = summoners[summonerNames[i]].ID
	}
	fmt.Println("Calculating...")
	scores, err := fantasy.CalculateScores(region, summonerIDs)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"scores": scores,
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
