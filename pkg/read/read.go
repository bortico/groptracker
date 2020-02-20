package read

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	structs "../structs"
)

func ParseArtists(url string) (structs.Artists, error) {
	var data structs.Artists
	res, err := http.Get(url)
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return data, err
	}
	json.Unmarshal(body, &data)
	return data, err
}

func ParseLocations(url string) (structs.Locations, error) {
	var data structs.Locations
	res, err := http.Get(url)
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return data, err
	}
	json.Unmarshal(body, &data)
	return data, err
}

func ParseDates(url string) (structs.Dates, error) {
	var data structs.Dates
	res, err := http.Get(url)
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return data, err
	}
	json.Unmarshal(body, &data)
	return data, err
}

func ParseRelation(url string) (structs.Relation, error) {
	var data structs.Relation
	res, err := http.Get(url)
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return data, err
	}
	json.Unmarshal(body, &data)
	return data, err
}
