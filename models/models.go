package models

type SiteData struct {
	UrlData  []UrlData `json:"data"`
	URLError error
}

type SiteDataResponse struct {
	UrlData []UrlData `json:"data"`
	Count   int       `json:"count"`
}

type UrlData struct {
	Url            string  `json:"url"`
	Views          int     `json:"views"`
	RelevanceScore float64 `json:"relevanceScore"`
}
