package webhelpers

import (
	"encoding/json"
	"log"
	"net/http"
)

//GetJSON makes a request and puts the json into target
func GetJSON(url string, target interface{}) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(target)

	//TODO: Remove this debug text
	log.Println("GET request to ", url, ":", target)

	return err
}
