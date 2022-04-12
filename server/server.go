package server

import (
	"assignment/httprequest"
	"assignment/models"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sort"
	"strconv"
	"sync"
)

var urls = []string{
	"https://raw.githubusercontent.com/assignment132/assignment/main/duckduckgo.json",
	"https://raw.githubusercontent.com/assignment132/assignment/main/google.json",
	"https://raw.githubusercontent.com/assignment132/assignment/main/wikipedia.json",
}

// GetData handles getData request and writes the response to ResponseWriter
func GetData(w http.ResponseWriter, req *http.Request) {
	key, limit, errCode, err := validateRequest(req)
	if err != nil {
		http.Error(w, err.Error(), errCode)
		log.Println(err)
		return
	}
	siteData := make(chan models.SiteData)

	var wg sync.WaitGroup
	for _, url := range urls {
		wg.Add(1)
		go httprequest.GetContent(url, siteData, &wg)
	}

	// close the channel in the background
	go func() {
		wg.Wait()
		close(siteData)
	}()

	var allSiteData models.SiteDataResponse
	errorInAPIs := false
	for data := range siteData {
		if data.URLError != nil {
			errorInAPIs = true
		}
		allSiteData.UrlData = append(allSiteData.UrlData, data.UrlData...)
	}
	allSiteData.Count = len(allSiteData.UrlData)

	if len(allSiteData.UrlData) == 0 && errorInAPIs {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	sortKey(allSiteData, key)
	if limit < len(allSiteData.UrlData) {
		allSiteData.UrlData = allSiteData.UrlData[0:limit]
		allSiteData.Count = limit
	}

	jsonResp, err := json.Marshal(allSiteData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("Error happened in JSON marshal: ", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
	log.Println("Response: ", allSiteData)
}

func validateRequest(req *http.Request) (string, int, int, error) {
	if req.Method != "GET" {
		return "", 0, http.StatusMethodNotAllowed, errors.New("method not allowed")
	}

	keys, ok := req.URL.Query()["sortKey"]
	if !ok || len(keys[0]) < 1 {
		return "", 0, http.StatusBadRequest, errors.New("url parameter 'sortKey' is missing")
	}
	key := keys[0]
	if key != "relevanceScore" && key != "views" {
		return "", 0, http.StatusBadRequest, errors.New("url parameter value for 'sortKey' is invalid")
	}

	limits, ok := req.URL.Query()["limit"]
	if !ok || len(limits[0]) < 1 {
		return "", 0, http.StatusBadRequest, errors.New("url parameter 'limit' is missing")
	}
	limit, err := strconv.Atoi(limits[0])
	if err != nil {
		return "", 0, http.StatusBadRequest, errors.New("Error while reading limit value: " + err.Error())
	}
	if limit < 1 || limit > 200 {
		return "", 0, http.StatusBadRequest, errors.New("url parameter value for 'limit' is invalid")
	}
	return key, limit, 0, nil
}

func sortKey(data models.SiteDataResponse, key string) {
	sort.Slice(data.UrlData, func(i, j int) bool {
		var result bool
		switch key {
		case "relevanceScore":
			result = data.UrlData[i].RelevanceScore < data.UrlData[j].RelevanceScore
		case "views":
			result = data.UrlData[i].Views < data.UrlData[j].Views
		}
		return result
	})
}
