package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
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

const (
	TL string = "┌" // top-left
	TM string = "┬" // top-mid
	TR string = "┐" // top-right
	ML string = "├" // mid-left
	MM string = "┼" // mid-mid
	MR string = "┤" // mid-right
	BL string = "└" // bottom-left
	BM string = "┴" // bottom-mid
	BR string = "┘" // bottom-right
	VE string = "─" // vertical
	HO string = "│" // horizontal
)

func adjustLeft(s string, l int) string {
	sLen := len(s)
	p := ""
	p = fmt.Sprint(s, strings.Repeat(" ", l-sLen))
	return p
}

func adjustRight(s string, l int) string {
	sLen := len(s)
	p := ""
	p = fmt.Sprint(strings.Repeat(" ", l-sLen), s)
	return p
}

type Table struct {
	Header         []string
	Rows           [][]string
	Size           int
	MaxColumnSizes []int
}

func NewTable(s int) *Table {
	t := &Table{Header: []string{}, Rows: [][]string{}, Size: s}
	t.MaxColumnSizes = make([]int, t.Size)
	return t
}

func (t *Table) SetHeader(columns []string) {
	if len(columns) != t.Size {
		log.Fatal("Number of columns does not match with size ", t.Size)
	}
	t.Header = columns
	t.UpdateMaxColumnSizes(columns)
}

func (t *Table) AppendRow(row []string) {
	if len(row) != t.Size {
		log.Fatal("Number of columns does not match with size ", t.Size)
	}
	t.Rows = append(t.Rows, row)
	t.UpdateMaxColumnSizes(row)
}

func (t *Table) UpdateMaxColumnSizes(records []string) {
	for i := 0; i < t.Size; i++ {
		if len(records[i]) > t.MaxColumnSizes[i] {
			t.MaxColumnSizes[i] = len(records[i])
		}
	}
}

func (t *Table) Display() {
	// Top border
	topBorder := ""
	topBorder = fmt.Sprint(topBorder, TL)
	for i, colSize := range t.MaxColumnSizes {
		topBorder = fmt.Sprint(topBorder, strings.Repeat(VE, colSize+2))
		if i != t.Size-1 {
			topBorder = fmt.Sprint(topBorder, TM)
		}
	}
	topBorder = fmt.Sprint(topBorder, TR)
	fmt.Println(topBorder)
	// Header
	header := ""
	header = fmt.Sprint(header, HO, " ")
	for i, col := range t.Header {
		header = fmt.Sprint(header, adjustLeft(col, t.MaxColumnSizes[i]))
		header = fmt.Sprint(header, " ", HO, " ")
	}
	fmt.Println(header)
	// Middle border
	middleBorder := ""
	middleBorder = fmt.Sprint(middleBorder, ML)
	for i, colSize := range t.MaxColumnSizes {
		middleBorder = fmt.Sprint(middleBorder, strings.Repeat(VE, colSize+2))
		if i != t.Size-1 {
			middleBorder = fmt.Sprint(middleBorder, MM)
		}
	}
	middleBorder = fmt.Sprint(middleBorder, MR)
	fmt.Println(middleBorder)
	// Rows
	for _, row := range t.Rows {
		rows := ""
		rows = fmt.Sprint(rows, HO, " ")
		for i, col := range row {
			if i == 0 {
				rows = fmt.Sprint(rows, adjustLeft(col, t.MaxColumnSizes[i]))
			} else {
				rows = fmt.Sprint(rows, adjustRight(col, t.MaxColumnSizes[i]))
			}
			rows = fmt.Sprint(rows, " ", HO, " ")
		}
		fmt.Println(rows)
	}
	// Bottom border
	bottomBorder := ""
	bottomBorder = fmt.Sprint(bottomBorder, BL)
	for i, colSize := range t.MaxColumnSizes {
		bottomBorder = fmt.Sprint(bottomBorder, strings.Repeat(VE, colSize+2))
		if i != t.Size-1 {
			bottomBorder = fmt.Sprint(bottomBorder, BM)
		}
	}
	bottomBorder = fmt.Sprint(bottomBorder, BR)
	fmt.Println(bottomBorder)
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

	// Generate table
	standings := NewTable(2)
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
