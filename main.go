package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"text/template"

	read "./pkg/read"
	structs "./pkg/structs"
)

var API structs.API

func main() {
	parseJSON()
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", indexHandle)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Listening server at port %v\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func parseJSON() {
	urlArtists := "https://groupietrackers.herokuapp.com/api/artists"
	urlLocations := "https://groupietrackers.herokuapp.com/api/locations"
	urlDates := "https://groupietrackers.herokuapp.com/api/dates"
	urlRelation := "https://groupietrackers.herokuapp.com/api/relation"
	API.Artists, _ = read.ParseArtists(urlArtists)
	API.Locations, _ = read.ParseLocations(urlLocations)
	API.Dates, _ = read.ParseDates(urlDates)
	API.Relation, _ = read.ParseRelation(urlRelation)
}

func indexHandle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		temp, err := template.ParseFiles("./static/templates/index.html")
		if err != nil {
			http.Error(w, "500 internal server error.", http.StatusInternalServerError)
			return
		}
		temp.Execute(w, API)
	case "POST":
		var toSearch string
		var searchType string

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "400 Bad request", http.StatusBadRequest)
			return
		}
		query, err := url.ParseQuery(string(body))
		if err != nil {
			http.Error(w, "400 Bad request", http.StatusBadRequest)
			return
		}

		for i, v := range query {
			switch i {
			case "toSearch":
				toSearch = v[0]
			case "searchType":
				searchType = v[0]
			default:
				http.Error(w, "400 Bad request", 400)
				return
			}
		}

		if searchType != "artist" && searchType != "member" && searchType != "location" && searchType != "firstAlbum" && searchType != "creationDate" {
			http.Error(w, "400 Bad request", http.StatusBadRequest)
			return
		}

		switch searchType {
		case "artist":
			sendArtist(w, r, toSearch)
		case "member":
			sendMember(w, r, toSearch)
		case "location":
			sendLocation(w, r, toSearch)
		case "firstAlbum":
			sendFirstAlbum(w, r, toSearch)
		case "creationDate":
			sendCreationDate(w, r, toSearch)
		}
	}
}

func sendArtist(w http.ResponseWriter, r *http.Request, toSearch string) {
	API.ID = -1
	for i := 0; i < 52; i++ {
		if strings.ToLower(API.Artists[i].Name) == strings.ToLower(toSearch) {
			API.ID = i
			break
		}
	}
	if API.ID == -1 {
		temp, err := template.ParseFiles("./static/templates/noresult.html")
		if err != nil {
			http.Error(w, "500 internal server error.", http.StatusInternalServerError)
			return
		}
		temp.Execute(w, toSearch)
	} else {
		temp, err := template.ParseFiles("./static/templates/artist.html")
		if err != nil {
			http.Error(w, "500 internal server error.", http.StatusInternalServerError)
			return
		}
		temp.Execute(w, API)
	}
}

func sendMember(w http.ResponseWriter, r *http.Request, toSearch string) {
	type MemberPage struct {
		Title  string
		Artist []string
	}
	var memberPage MemberPage
	for i := 0; i < 52; i++ {
		for _, member := range API.Artists[i].Members {
			if strings.ToLower(member) == strings.ToLower(toSearch) {
				memberPage.Title = member + "<br>is a member of"
				memberPage.Artist = append(memberPage.Artist, API.Artists[i].Name)
			}
		}
	}
	if memberPage.Title == "" {
		temp, err := template.ParseFiles("./static/templates/noresult.html")
		if err != nil {
			http.Error(w, "500 internal server error.", http.StatusInternalServerError)
			return
		}
		temp.Execute(w, toSearch)
	} else {
		temp, err := template.ParseFiles("./static/templates/member.html")
		if err != nil {
			http.Error(w, "500 internal server error.", http.StatusInternalServerError)
			return
		}
		temp.Execute(w, memberPage)
	}
}

func sendLocation(w http.ResponseWriter, r *http.Request, toSearch string) {
	type LocationPage struct {
		Title   string
		Artists []string
	}
	var locationPage LocationPage

	for i, all := range API.Locations.Index {
		for _, location := range all.Locations {
			if strings.ToLower(location) == strings.ToLower(toSearch) {
				locationPage.Title = "Concerts in " + location
				locationPage.Artists = append(locationPage.Artists, API.Artists[i].Name)
				break
			}
		}
	}
	if locationPage.Title == "" {
		temp, err := template.ParseFiles("./static/templates/noresult.html")
		if err != nil {
			http.Error(w, "500 internal server error.", http.StatusInternalServerError)
			return
		}
		temp.Execute(w, toSearch)
	} else {
		temp, err := template.ParseFiles("./static/templates/location.html")
		if err != nil {
			http.Error(w, "500 internal server error.", http.StatusInternalServerError)
			return
		}
		temp.Execute(w, locationPage)
	}
}

func sendFirstAlbum(w http.ResponseWriter, r *http.Request, toSearch string) {
	type FirstAlbumPage struct {
		Title   string
		Artists []string
	}
	var firstAlbumPage FirstAlbumPage

	for i, artist := range API.Artists {
		if strings.ToLower(artist.FirstAlbum) == strings.ToLower(toSearch) {
			firstAlbumPage.Title = "Artists / Bands relased their first album in " + artist.FirstAlbum
			firstAlbumPage.Artists = append(firstAlbumPage.Artists, API.Artists[i].Name)
		}
	}
	if firstAlbumPage.Title == "" {
		temp, err := template.ParseFiles("./static/templates/noresult.html")
		if err != nil {
			http.Error(w, "500 internal server error.", http.StatusInternalServerError)
			return
		}
		temp.Execute(w, toSearch)
	} else {
		temp, err := template.ParseFiles("./static/templates/firstalbum.html")
		if err != nil {
			http.Error(w, "500 internal server error.", http.StatusInternalServerError)
			return
		}
		temp.Execute(w, firstAlbumPage)
	}
}

func sendCreationDate(w http.ResponseWriter, r *http.Request, toSearch string) {
	year, _ := strconv.Atoi(toSearch)
	type CreationDatePage struct {
		Title   string
		Artists []string
	}
	var creationDatePage CreationDatePage

	for i, artist := range API.Artists {
		if artist.CreationDate == year {
			creationDatePage.Title = "Artists / Bands created in " + toSearch
			creationDatePage.Artists = append(creationDatePage.Artists, API.Artists[i].Name)
		}
	}
	if creationDatePage.Title == "" {
		temp, err := template.ParseFiles("./static/templates/noresult.html")
		if err != nil {
			http.Error(w, "500 internal server error.", http.StatusInternalServerError)
			return
		}
		temp.Execute(w, toSearch)
	} else {
		temp, err := template.ParseFiles("./static/templates/creationdate.html")
		if err != nil {
			http.Error(w, "500 internal server error.", http.StatusInternalServerError)
			return
		}
		temp.Execute(w, creationDatePage)
	}
}
