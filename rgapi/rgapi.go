package rgapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	BaseURL string = "api.riotgames.com/lol"

	apiKey string = ""

	//ErrAPIKeyNotSet is the error returned when no global API key has been set
	ErrAPIKeyNotSet = errors.New("rgapi: API key has not been set. If you need a key visit https://developer.riotgames.com/")
	shortRateChan   rateChan
	longRateChan    rateChan
)

type rateChan struct {
	RateQueue   chan bool
	TriggerChan chan bool
}

//HttpError contains the http status code of the erro
type HttpError struct {
	StatusCode int
}

func SetAPIKey(key string) {
	apiKey = key
}

func isKeySet() bool {
	return apiKey != ""
}

func apiGet(region, requestURL string, target interface{}) error {
	if !isKeySet() {
		return ErrAPIKeyNotSet
	}
	checkRateLimiter(shortRateChan)
	checkRateLimiter(longRateChan)
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://%s.%s%s", region, BaseURL, requestURL), nil)
	req.Header.Add("X-Riot-Token", apiKey)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	checkTimeTrigger(shortRateChan)
	checkTimeTrigger(longRateChan)
	if resp.StatusCode != http.StatusOK {
		return HttpError{StatusCode: resp.StatusCode}
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, target)
	if err != nil {
		return err
	}
	return err
}

//SetShortRateLimit allows a custom rate limit to be set. For at the time of this writing the default
//for a development API key is 10 requests every 10 seconds
func SetShortRateLimit(numrequests int, pertime time.Duration) {
	shortRateChan = rateChan{
		RateQueue:   make(chan bool, numrequests),
		TriggerChan: make(chan bool),
	}
	go rateLimitHandler(shortRateChan, pertime)
}

func rateLimitHandler(RateChan rateChan, pertime time.Duration) {
	returnChan := make(chan bool)
	go timeTriggerWatcher(RateChan.TriggerChan, returnChan)
	for {
		<-returnChan
		<-time.After(pertime)
		go timeTriggerWatcher(RateChan.TriggerChan, returnChan)
		length := len(RateChan.RateQueue)
		for i := 0; i < length; i++ {
			<-RateChan.RateQueue
		}
	}
}

func timeTriggerWatcher(timeTrigger chan bool, returnChan chan bool) {
	timeTrigger <- true
	returnChan <- true
}

//SetLongRateLimit allows a custom rate limit to be set. For at the time of this writing the default
//for a development API key is 500 requests every 10 minutes
func SetLongRateLimit(numrequests int, pertime time.Duration) {
	longRateChan = rateChan{
		RateQueue:   make(chan bool, numrequests),
		TriggerChan: make(chan bool),
	}
	go rateLimitHandler(longRateChan, pertime)
}

func checkRateLimiter(RateChan rateChan) {
	if RateChan.RateQueue != nil && RateChan.TriggerChan != nil {
		RateChan.RateQueue <- true
	}
}

func checkTimeTrigger(RateChan rateChan) {
	if RateChan.RateQueue != nil && RateChan.TriggerChan != nil {
		select {
		case <-RateChan.TriggerChan:
		default:
		}
	}
}

//Error prints the error message for a HttpError
func (err HttpError) Error() string {
	return fmt.Sprintf("Error: HTTP Status %d", err.StatusCode)
}

// The entire rgapi package is modified from: https://github.com/TrevorSStone/goriot
//
// The MIT License (MIT)
//
// Copyright (c) 2013 TrevorSStone
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
