package main

import (
	"log"
	"os"

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

	controller.Route()
}
