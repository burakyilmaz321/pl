package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
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
	Header           []string
	Rows             [][]string
	Size             int
	MaxColumnSizes   []int
	Padding          int
	ColumnAlignments []int
}

func NewTable(s int, alignments []int) *Table {
	t := &Table{Header: []string{}, Rows: [][]string{}, Size: s}
	t.MaxColumnSizes = make([]int, t.Size)
	t.Padding = 1
	t.ColumnAlignments = alignments
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

func (t *Table) BuildBorder(left string, middle string, right string) string {
	cols := make([]string, t.Size)
	for i := 0; i < t.Size; i++ {
		cols[i] = fmt.Sprint(strings.Repeat(VE, t.MaxColumnSizes[i]+t.Padding*2))
	}
	return fmt.Sprint(left, strings.Join(cols, middle), right)
}

func (t *Table) BuildRowLine(row []string, columnAlignments []int) string {
	cols := make([]string, t.Size)
	for i := 0; i < t.Size; i++ {
		alignment := columnAlignments[i]
		switch alignment {
		case 0:
			cols[i] = adjustLeft(row[i], t.MaxColumnSizes[i])
		case 1:
			cols[i] = adjustRight(row[i], t.MaxColumnSizes[i])
		}
		cols[i] = fmt.Sprint(strings.Repeat(" ", t.Padding), cols[i], strings.Repeat(" ", t.Padding))
	}
	return fmt.Sprint(HO, strings.Join(cols, HO), HO)
}

func (t *Table) Display() {
	// Top border
	topBorder := t.BuildBorder(TL, TM, TR)
	fmt.Println(topBorder)
	// Header
	header := t.BuildRowLine(t.Header, []int{0, 0})
	fmt.Println(header)
	// Middle border
	middleBorder := t.BuildBorder(ML, MM, MR)
	fmt.Println(middleBorder)
	// Rows
	for _, row := range t.Rows {
		fmt.Println(t.BuildRowLine(row, t.ColumnAlignments))
	}
	// Bottom border
	bottomBorder := t.BuildBorder(BL, BM, BR)
	fmt.Println(bottomBorder)
}

func Get(url string, params map[string]string) (*http.Response, error) {
	// Create new http client
	client := &http.Client{}

	// Create new http request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Add query string params to the request
	queryParams := req.URL.Query()

	for param, value := range params {
		queryParams.Add(param, value)
	}

	req.URL.RawQuery = queryParams.Encode()

	// Execute the request
	res, err := client.Do(req)

	return res, err
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
	res, err := Get(HOSTNAME, queryParams)

	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	// Parse response as struct
	var dat Response
	json.NewDecoder(res.Body).Decode(&dat)

	// Generate table
	standings := NewTable(2, []int{0, 1})
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
