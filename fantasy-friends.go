package main

import (
	"log"
	"net/http"
	"os"

	"github.com/TrevorSStone/goriot"
	"github.com/gin-gonic/gin"
)

func main() {
	//Get port number to listen for
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	//Set the Riot API key
	riotAPIKey := os.Getenv("RIOT_API_KEY")
	if riotAPIKey == "" {
		log.Fatal("$RIOT_API_KEY must be set")
	}
	goriot.SetAPIKey(riotAPIKey)

	//Set up routing
	router := gin.Default()

	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	//TODO: Organize routing and remove this test code =============================================================
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	router.POST("/summoner", func(c *gin.Context) {
		summonerName := goriot.NormalizeSummonerName(c.PostForm("summonerName"))[0]

		//Attempt to get the summoner ID
		summoners, err := goriot.SummonerByName("na", summonerName)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "failure",
				"error":  err.Error(),
			})
		} else {
			if summoner, ok := summoners[summonerName]; ok {
				c.JSON(http.StatusOK, gin.H{
					"status":       "success",
					"summonerID":   summoner.ID,
					"summonerName": summoner.Name,
				})
			} else {
				c.JSON(http.StatusBadRequest, gin.H{
					"status": "failure",
					"error":  "Summoner not found.",
				})
			}
		}

	})
	//==============================================================================================================

	router.Run(":" + port)
}
