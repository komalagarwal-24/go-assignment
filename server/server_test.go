package server

import (
	"assignment/httprequest"
	"assignment/models"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetData(t *testing.T) {
	httprequest.GetContent = func(url string, siteData chan models.SiteData, wg *sync.WaitGroup) {
		defer wg.Done()
		switch url {
		case urls[0]:
			siteData <- models.SiteData{
				UrlData: []models.UrlData{
					{
						Url:            "www.yahoo.com/abc6",
						Views:          6000,
						RelevanceScore: 0.6,
					},
					{
						Url:            "www.yahoo.com/abc7",
						Views:          7000,
						RelevanceScore: 0.7,
					},
					{
						Url:            "www.yahoo.com/abc8",
						Views:          8000,
						RelevanceScore: 0.8,
					},
				},
			}
		case urls[1]:
			siteData <- models.SiteData{
				UrlData: []models.UrlData{
					{
						Url:            "www.example.com/abc1",
						Views:          1000,
						RelevanceScore: 0.4,
					},
					{
						Url:            "www.example.com/abc2",
						Views:          2000,
						RelevanceScore: 0.5,
					},
					{
						Url:            "www.example.com/abc3",
						Views:          3000,
						RelevanceScore: 0.65,
					},
				},
			}
		case urls[2]:
			siteData <- models.SiteData{
				UrlData: []models.UrlData{
					{
						Url:            "www.wikipedia.com/abc1",
						Views:          11000,
						RelevanceScore: 0.1,
					},
					{
						Url:            "www.wikipedia.com/abc2",
						Views:          12000,
						RelevanceScore: 0.2,
					},
					{
						Url:            "www.wikipedia.com/abc3",
						Views:          13000,
						RelevanceScore: 0.3,
					},
				},
			}
		}
	}

	// Create a request to pass to handler
	req, err := http.NewRequest("GET", "/getData?sortKey=views&limit=5", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetData)

	// Call ServeHTTP method directly and pass in Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("GetDataSortOnViews: handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := models.SiteDataResponse{
		UrlData: []models.UrlData{
			{
				Url:            "www.example.com/abc1",
				Views:          1000,
				RelevanceScore: 0.4,
			},
			{
				Url:            "www.example.com/abc2",
				Views:          2000,
				RelevanceScore: 0.5,
			},
			{
				Url:            "www.example.com/abc3",
				Views:          3000,
				RelevanceScore: 0.65,
			},
			{
				Url:            "www.yahoo.com/abc6",
				Views:          6000,
				RelevanceScore: 0.6,
			},
			{
				Url:            "www.yahoo.com/abc7",
				Views:          7000,
				RelevanceScore: 0.7,
			},
		},
		Count: 5,
	}

	var data models.SiteDataResponse
	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		t.Fatal(err)
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		t.Fatal(err)
	}

	if !equalData(data, expected) {
		t.Errorf("GetDataSortOnViews: handler returned unexpected body: got %v want %v",
			data, expected)
	}

	req, err = http.NewRequest("GET", "/getData?sortKey=relevanceScore&limit=6", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("GetDataSortOnRelevanceScore: handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected = models.SiteDataResponse{
		UrlData: []models.UrlData{
			{
				Url:            "www.wikipedia.com/abc1",
				Views:          11000,
				RelevanceScore: 0.1,
			},
			{
				Url:            "www.wikipedia.com/abc2",
				Views:          12000,
				RelevanceScore: 0.2,
			},
			{
				Url:            "www.wikipedia.com/abc3",
				Views:          13000,
				RelevanceScore: 0.3,
			},
			{
				Url:            "www.example.com/abc1",
				Views:          1000,
				RelevanceScore: 0.4,
			},
			{
				Url:            "www.example.com/abc2",
				Views:          2000,
				RelevanceScore: 0.5,
			},
			{
				Url:            "www.yahoo.com/abc6",
				Views:          6000,
				RelevanceScore: 0.6,
			},
		},
		Count: 6,
	}

	body, err = ioutil.ReadAll(rr.Body)
	if err != nil {
		t.Fatal(err)
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		t.Fatal(err)
	}

	if !equalData(data, expected) {
		t.Errorf("GetDataSortOnRelevanceScore: handler returned unexpected body: got %v want %v",
			data, expected)
	}

}

func TestGetDataReturnsNoData(t *testing.T) {
	httprequest.GetContent = func(url string, siteData chan models.SiteData, wg *sync.WaitGroup) {
		defer wg.Done()
		switch url {
		case urls[0]:
			siteData <- models.SiteData{}
		case urls[1]:
			siteData <- models.SiteData{}
		case urls[2]:
			siteData <- models.SiteData{}
		}
	}

	req, err := http.NewRequest("GET", "/getData?sortKey=views&limit=5", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetData)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("TestGetDataReturnsNoData: handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestGetDataReturnsDataWhenErrorInSomeAPIs(t *testing.T) {
	httprequest.GetContent = func(url string, siteData chan models.SiteData, wg *sync.WaitGroup) {
		defer wg.Done()
		switch url {
		case urls[0]:
			siteData <- models.SiteData{
				URLError: errors.New("Internal error"),
			}
		case urls[1]:
			siteData <- models.SiteData{
				UrlData: []models.UrlData{
					{
						Url:            "www.wikipedia.com/abc1",
						Views:          11000,
						RelevanceScore: 0.1,
					},
					{
						Url:            "www.wikipedia.com/abc2",
						Views:          12000,
						RelevanceScore: 0.2,
					},
				},
			}
		case urls[2]:
			siteData <- models.SiteData{
				UrlData: []models.UrlData{
					{
						Url:            "www.example.com/abc1",
						Views:          1000,
						RelevanceScore: 0.4,
					},
					{
						Url:            "www.example.com/abc2",
						Views:          2000,
						RelevanceScore: 0.5,
					},
				},
			}
		}
	}

	req, err := http.NewRequest("GET", "/getData?sortKey=views&limit=50", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetData)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("TestGetDataReturnsDataWhenErrorInSomeAPIs: handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := models.SiteDataResponse{
		UrlData: []models.UrlData{
			{
				Url:            "www.example.com/abc1",
				Views:          1000,
				RelevanceScore: 0.4,
			},
			{
				Url:            "www.example.com/abc2",
				Views:          2000,
				RelevanceScore: 0.5,
			},
			{
				Url:            "www.wikipedia.com/abc1",
				Views:          11000,
				RelevanceScore: 0.1,
			},
			{
				Url:            "www.wikipedia.com/abc2",
				Views:          12000,
				RelevanceScore: 0.2,
			},
		},
		Count: 4,
	}

	var data models.SiteDataResponse
	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		t.Fatal(err)
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		t.Fatal(err)
	}

	if !equalData(data, expected) {
		t.Errorf("TestGetDataReturnsDataWhenErrorInSomeAPIs: handler returned unexpected body: got %v want %v",
			data, expected)
	}
}

func TestGetDataInvalidRequest(t *testing.T) {
	req, err := http.NewRequest("POST", "/getData?sortKey=views&limit=5", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetData)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("TestGetDataInvalidRequest: handler returned wrong status code: got %v want %v",
			status, http.StatusMethodNotAllowed)
	}
}

func TestGetDataInvalidRequestParams1(t *testing.T) {
	req, err := http.NewRequest("GET", "/getData?sortKey=views&limit=500", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetData)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("TestGetDataInvalidRequestParams1: handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestGetDataInvalidRequestParams2(t *testing.T) {
	req, err := http.NewRequest("GET", "/getData?sortKey=abc&limit=50", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetData)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("TestGetDataInvalidRequestParams2: handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestGetDataInvalidRequestParams3(t *testing.T) {
	req, err := http.NewRequest("GET", "/getData?sortKey=views", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetData)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("TestGetDataInvalidRequestParams3: handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestGetDataReturnsError(t *testing.T) {
	httprequest.GetContent = func(url string, siteData chan models.SiteData, wg *sync.WaitGroup) {
		defer wg.Done()
		switch url {
		case urls[0]:
			siteData <- models.SiteData{}
		case urls[1]:
			siteData <- models.SiteData{
				URLError: errors.New("Some error"),
			}
		case urls[2]:
			siteData <- models.SiteData{}
		}
	}

	req, err := http.NewRequest("GET", "/getData?sortKey=views&limit=5", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetData)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("TestGetDataReturnsError: handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}

func Test_validateRequest(t *testing.T) {
	type args struct {
		req *http.Request
	}
	tests := []struct {
		name    string
		args    args
		key     string
		limit   int
		errCode int
		wantErr bool
	}{
		{
			name: "TestSuccessCase",
			args: args{
				req: &http.Request{
					Method: "GET",
					URL: &url.URL{
						RawQuery: "sortKey=views&limit=5",
					},
				},
			},
			key:     "views",
			limit:   5,
			errCode: 0,
			wantErr: false,
		},
		{
			name: "TestInvalidMethod",
			args: args{
				req: &http.Request{
					Method: "POST",
					URL: &url.URL{
						RawQuery: "sortKey=views&limit=5",
					},
				},
			},
			key:     "",
			limit:   0,
			errCode: http.StatusMethodNotAllowed,
			wantErr: true,
		},
		{
			name: "TestSortKeyMissing",
			args: args{
				req: &http.Request{
					Method: "GET",
					URL: &url.URL{
						RawQuery: "limit=5",
					},
				},
			},
			key:     "",
			limit:   0,
			errCode: http.StatusBadRequest,
			wantErr: true,
		},
		{
			name: "TestSortKeyInvalidValue",
			args: args{
				req: &http.Request{
					Method: "GET",
					URL: &url.URL{
						RawQuery: "sortKey=abc&limit=5",
					},
				},
			},
			key:     "",
			limit:   0,
			errCode: http.StatusBadRequest,
			wantErr: true,
		},
		{
			name: "TestLimitMissing",
			args: args{
				req: &http.Request{
					Method: "GET",
					URL: &url.URL{
						RawQuery: "sortKey=views",
					},
				},
			},
			key:     "",
			limit:   0,
			errCode: http.StatusBadRequest,
			wantErr: true,
		},
		{
			name: "TestLimitInvalidValue",
			args: args{
				req: &http.Request{
					Method: "GET",
					URL: &url.URL{
						RawQuery: "sortKey=views&limit=500",
					},
				},
			},
			key:     "",
			limit:   0,
			errCode: http.StatusBadRequest,
			wantErr: true,
		},
		{
			name: "TestLimitInvalidValue",
			args: args{
				req: &http.Request{
					Method: "GET",
					URL: &url.URL{
						RawQuery: "sortKey=views&limit=abc",
					},
				},
			},
			key:     "",
			limit:   0,
			errCode: http.StatusBadRequest,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2, err := validateRequest(tt.args.req)
			assert.Equal(t, tt.wantErr, (err != nil))
			assert.Equal(t, tt.key, got)
			assert.Equal(t, tt.limit, got1)
			assert.Equal(t, tt.errCode, got2)
		})
	}
}

func Test_sortKey(t *testing.T) {
	type args struct {
		data models.SiteDataResponse
		key  string
	}
	tests := []struct {
		name         string
		args         args
		expectedData models.SiteDataResponse
	}{
		{
			name: "TestSortRelevanceScore",
			args: args{
				data: models.SiteDataResponse{
					UrlData: []models.UrlData{
						{
							Url:            "www.wikipedia.com/abc1",
							Views:          100,
							RelevanceScore: 0.5,
						},
						{
							Url:            "www.wikipedia.com/abc2",
							Views:          200,
							RelevanceScore: 0.4,
						},
						{
							Url:            "www.wikipedia.com/abc3",
							Views:          300,
							RelevanceScore: 0.3,
						},
					},
					Count: 3,
				},
				key: "relevanceScore",
			},
			expectedData: models.SiteDataResponse{
				UrlData: []models.UrlData{
					{
						Url:            "www.wikipedia.com/abc3",
						Views:          300,
						RelevanceScore: 0.3,
					},
					{
						Url:            "www.wikipedia.com/abc2",
						Views:          200,
						RelevanceScore: 0.4,
					},
					{
						Url:            "www.wikipedia.com/abc1",
						Views:          100,
						RelevanceScore: 0.5,
					},
				},
				Count: 3,
			},
		},
		{
			name: "TestSortViews",
			args: args{
				data: models.SiteDataResponse{
					UrlData: []models.UrlData{
						{
							Url:            "www.wikipedia.com/abc1",
							Views:          300,
							RelevanceScore: 0.3,
						},
						{
							Url:            "www.wikipedia.com/abc2",
							Views:          200,
							RelevanceScore: 0.4,
						},
						{
							Url:            "www.wikipedia.com/abc3",
							Views:          100,
							RelevanceScore: 0.5,
						},
					},
					Count: 3,
				},
				key: "views",
			},
			expectedData: models.SiteDataResponse{
				UrlData: []models.UrlData{
					{
						Url:            "www.wikipedia.com/abc3",
						Views:          100,
						RelevanceScore: 0.5,
					},
					{
						Url:            "www.wikipedia.com/abc2",
						Views:          200,
						RelevanceScore: 0.4,
					},
					{
						Url:            "www.wikipedia.com/abc1",
						Views:          300,
						RelevanceScore: 0.3,
					},
				},
				Count: 3,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sortKey(tt.args.data, tt.args.key)
			if !equalData(tt.args.data, tt.expectedData) {
				t.Errorf("sortKey() data = %v, want %v", tt.args.data, tt.expectedData)
			}
		})
	}
}

func equalData(data1 models.SiteDataResponse, data2 models.SiteDataResponse) bool {
	if len(data1.UrlData) != len(data2.UrlData) {
		return false
	}
	for data1.Count != data2.Count {
		return false
	}
	for i, data := range data1.UrlData {
		if data.RelevanceScore != data2.UrlData[i].RelevanceScore || data.Views != data2.UrlData[i].Views {
			return false
		}
	}
	return true
}
