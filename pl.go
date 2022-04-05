package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const hostname string = "https://site.web.api.espn.com/apis/v2/sports/soccer/eng.1/standings"

type Response struct {
	Uid         string `json:"uid"`
	Id          string `json:"id"`
	Name        string `json:"name"`
	Abbrevation string `json:"abbrevation"`
	Children    []struct {
		Uid         string `json:"uid"`
		Id          string `json:"id"`
		Name        string `json:"name"`
		Abbrevation string `json:"abbrevation"`
		Standings   struct {
			Id         string        `json:"id"`
			Name       string        `json:"name"`
			Links      []interface{} `json:"links"`
			Season     string        `json:"season"`
			SeasonType string        `json:"seasonType"`
			Entries    []struct {
				Team struct {
					Id               string        `json:"id"`
					Uid              string        `json:"uid"`
					Location         string        `json:"location"`
					Name             string        `json:"name"`
					Abbreviation     string        `json:"abbreviation"`
					DisplayName      string        `json:"displayName"`
					ShortDisplayName string        `json:"shortDisplayName"`
					IsActive         string        `json:"isActive"`
					Logos            []interface{} `json:"logos"`
					Links            []interface{} `json:"links"`
				} `json:"team"`
				Note struct {
					Color       string `json:"color"`
					Description string `json:"description"`
					Rank        string `json:"rank"`
				} `json:"note"`
				Stats []struct {
					Name             string `json:"name"`
					DisplayName      string `json:"displayName"`
					ShortDisplayName string `json:"shortDisplayName"`
					Description      string `json:"description"`
					Abbreviation     string `json:"abbreviation"`
					Type             string `json:"type"`
					Value            string `json:"value"`
					DisplayValue     string `json:"displayValue"`
				} `json:"stats"`
			} `json:"entries"`
		} `json:"standings"`
	} `json:"children"`
	Seasons []interface{} `json:"seasons"`
}

func main() {
	// Create new http client
	client := &http.Client{}

	// Create new http request
	req, err := http.NewRequest("GET", hostname, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Add query string params to the request
	queryParams := req.URL.Query()
	queryParams.Add("region", "us")
	queryParams.Add("lang", "en")
	queryParams.Add("contentorigin", "soccernet")
	queryParams.Add("season", "2021")
	queryParams.Add("sort", "rank")
	req.URL.RawQuery = queryParams.Encode()

	// Execute the request
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	// Parse response as struct
	var dat Response
	json.NewDecoder(res.Body).Decode(&dat)

	// Print table
	fmt.Println("+------|--------+")
	fmt.Println("| Team | Points |")
	fmt.Println("+======|========+")
	for _, i := range dat.Children[0].Standings.Entries {
		var points string
		for _, j := range i.Stats {
			if j.Name == "points" {
				points = j.DisplayValue
			}
		}
		fmt.Println("|", i.Team.Abbreviation, " |    ", points, "|")
		fmt.Println("+------|--------+")
	}
}
