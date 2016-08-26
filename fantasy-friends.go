package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/xelaadryth/fantasy-friends/riot"
)

type privateDatatype struct {
	RiotAPIKey string
}

const privateDataFilename = "private.json"

var privateData = privateDatatype{}

func importPrivateData() error {
	f, err := os.Open(privateDataFilename)
	if err != nil {
		return err
	}
	return json.NewDecoder(f).Decode(&privateData)
}

func main() {
	//Import any settings that should be hidden such as passwords and keys
	err := importPrivateData()
	if err != nil {
		log.Fatal("Unable to import private data from private.json; core functionality will not work.")
	}
	riot.SetAPIKey(privateData.RiotAPIKey)

	//Start the web server
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.Default()

	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	//TODO: Organize routing and remove this test code ===================================================================
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	router.POST("/summoner", func(c *gin.Context) {
		summonerName := c.PostForm("summonerName")

		//Attempt to get the summoner ID
		summonerID, err := riot.GetSummonerID("na", summonerName)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":       "failure",
				"error":        err.Error(),
				"summonerName": summonerName,
				"summonerID":   summonerID,
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"status":       "success",
			"summonerName": summonerName,
			"summonerID":   summonerID,
		})

	})
	//====================================================================================================================

	router.Run(":" + port)
}
