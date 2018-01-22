package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/xelaadryth/fantasy-friends/controller"
	"github.com/xelaadryth/fantasy-friends/database"
	"github.com/xelaadryth/fantasy-friends/rgapi"
)

func main() {
	//Set the Riot API key
	riotAPIKey := os.Getenv("RIOT_API_KEY")
	if riotAPIKey == "" {
		log.Fatal("$RIOT_API_KEY must be set")
	}
	rgapi.SetAPIKey(riotAPIKey)

	// TODO: Make this configurable live
	var shortLimit, longLimit int
	var err error
	shortLimit, err = strconv.Atoi(os.Getenv("SHORT_RATE_LIMIT"))
	if err != nil {
		log.Fatal("Unable to get the short rate limit")
	} else {
		rgapi.SetShortRateLimit(shortLimit, 10*time.Second)
	}
	longLimit, err = strconv.Atoi(os.Getenv("LONG_RATE_LIMIT"))
	if err != nil {
		log.Fatal("Unable to get the long rate limit")
	} else {
		rgapi.SetLongRateLimit(longLimit, 10*time.Minute)
	}

	//Set up DB
	err = database.Connect()
	if err != nil {
		log.Fatal(err.Error())
	}

	database.PreparePepper()

	controller.Route()
}
