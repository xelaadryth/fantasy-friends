package main

import (
	"log"
	"os"
	"time"

	"github.com/TrevorSStone/goriot"
	"github.com/xelaadryth/fantasy-friends/controller"
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

	controller.Route()
}
