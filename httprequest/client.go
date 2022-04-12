package httprequest

import (
	"assignment/models"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

const (
	retries = 3 // default max retry count
)

type API struct {
	Client  *http.Client
	BaseURL string
}

var GetContent = getContent

// getContent executes GET request with url and writes response on channel
func getContent(url string, siteData chan models.SiteData, wg *sync.WaitGroup) {
	defer wg.Done()
	client := http.Client{
		Timeout: 2 * time.Second,
	}

	api := API{
		Client:  &client,
		BaseURL: url,
	}
	data := api.ExecuteAPI()
	siteData <- data
}

// ExecuteAPI makes get request to the api and returns response
func (api *API) ExecuteAPI() models.SiteData {
	var data models.SiteData
	sleep := 2 * time.Second

	for i := 0; i < retries; i++ {
		if i > 0 {
			log.Println("Retrying after error: ", api.BaseURL, data.URLError)
			time.Sleep(sleep)
			sleep *= 2
		}
		resp, err := api.Client.Get(api.BaseURL)
		if err != nil {
			log.Println("Error while making http request: ", api.BaseURL, err)
			data.URLError = err
			continue
		}
		data, err = parseResponse(api.BaseURL, resp)
		if err != nil {
			data.URLError = err
			continue
		}

		return data
	}

	return data
}

func parseResponse(url string, resp *http.Response) (models.SiteData, error) {
	defer resp.Body.Close()
	var data models.SiteData

	if resp.StatusCode != http.StatusOK {
		log.Println("Status code is not 200 for: " + url)
		return data, errors.New("Status code is not 200 for: " + url)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error while reading http response: ", url, err)
		return data, err
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Println("Error while parsing http response: ", url, err)
		return data, err
	}
	return data, nil
}
