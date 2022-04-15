package main

import (
	"encoding/json"
	"log"

	"github.com/burakyilmaz321/pl/pkg/requests"
	"github.com/burakyilmaz321/pl/pkg/table"
)

const HOSTNAME string = "https://site.web.api.espn.com/apis/v2/sports/soccer/eng.1/standings"

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
	// Make GET request
	queryParams := map[string]string{
		"region":        "us",
		"lang":          "en",
		"contentorigin": "soccernet",
		"season":        "2021",
		"sort":          "rank",
	}
	res, err := requests.Get(HOSTNAME, queryParams)

	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	// Parse response as struct
	var dat Response
	json.NewDecoder(res.Body).Decode(&dat)

	// Generate table
	standings := table.NewTable(2, []int{0, 1})
	standings.SetHeader([]string{"Team", "Points"})
	for _, i := range dat.Children[0].Standings.Entries {
		var points string
		for _, j := range i.Stats {
			if j.Name == "points" {
				points = j.DisplayValue
			}
		}
		standings.AppendRow([]string{i.Team.DisplayName, points})
	}
	// Print table
	standings.Display()
}
