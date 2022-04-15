package requests

import (
	"log"
	"net/http"
)

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
