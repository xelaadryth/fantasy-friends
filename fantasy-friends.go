package main

import (
	"log"
	"os"
	"time"

	"github.com/TrevorSStone/goriot"
	"github.com/xelaadryth/fantasy-friends/controller"
	"github.com/xelaadryth/fantasy-friends/database"
)

func main() {
	//Set the Riot API key
	riotAPIKey := os.Getenv("RIOT_API_KEY")
	if riotAPIKey == "" {
		log.Fatal("$RIOT_API_KEY must be set")
	}
	goriot.SetAPIKey(riotAPIKey)
	goriot.SetSmallRateLimit(10, 10*time.Second)
	goriot.SetLongRateLimit(500, 10*time.Minute)

	//Set up DB
	err := database.Connect()
	if err != nil {
		log.Fatal(err.Error())
	}

	controller.Route()
}
