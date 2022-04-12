package httprequest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecuteAPI(t *testing.T) {
	// Start a local HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(`{
			"data": [
			  {
				"url": "www.wikipedia.com/abc2",
				"views": 12000,
				"relevanceScore": 0.2
			  },
			  {
				"url": "www.wikipedia.com/abc4",
				"views": 13000,
				"relevanceScore": 0.4
			  },
			  {
				"url": "www.wikipedia.com/abc3",
				"views": 14000,
				"relevanceScore": 0.3
			  },
			  {
				"url": "www.wikipedia.com/abc5",
				"views": 15000,
				"relevanceScore": 0.5
			  },
			  {
				"url": "www.wikipedia.com/abc1",
				"views": 11000,
				"relevanceScore": 0.1
			  }
			]
		  }`))
	}))
	defer server.Close()

	api := API{server.Client(), server.URL}
	data := api.ExecuteAPI()

	assert.Equal(t, len(data.UrlData), 5)
}

func TestExecuteAPIRetry(t *testing.T) {
	var retryCount = 0
	// Start a local HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		retryCount++
		rw.Write([]byte(`invalid response`))
	}))
	defer server.Close()

	api := API{server.Client(), server.URL}
	data := api.ExecuteAPI()

	assert.Equal(t, retryCount, retries)
	assert.Equal(t, len(data.UrlData), 0)
}

func TestExecuteAPIRetryTwo(t *testing.T) {
	var retryCount = 0
	// Start a local HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		retryCount++
		rw.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	api := API{server.Client(), server.URL}
	data := api.ExecuteAPI()

	assert.Equal(t, retryCount, retries)
	assert.Equal(t, len(data.UrlData), 0)
}
