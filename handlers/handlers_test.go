package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"paragliding/database"
	"testing"
)

//Testing the Redirect func
func TestRedirect(t *testing.T) {
	TestTable := []struct {
		method string
		url    string
		code   int
	}{
		{method: "GET", url: "http://localhost:8080/paragliding/", code: 301},
		{method: "GET", url: "http://localhost:8080/paragliding/rr", code: 404},
		{method: "POST", url: "http://localhost:8080/paragliding/", code: 404},
	}
	for _, testCase := range TestTable {
		req, err := http.NewRequest(testCase.method, testCase.url, nil)
		if err != nil {
			t.Errorf("Unexpected error, %d", err)
		}

		writer := httptest.NewRecorder()
		Redirect(writer, req)
		if writer.Code != testCase.code {
			t.Errorf("got resp code %v, expected %v", writer.Code, testCase.code)
		}
	}
}

// Testing the function Index
func TestIndex(t *testing.T) {
	TestTable := []struct {
		method string
		url    string
		code   int
	}{
		{method: "GET", url: "http://localhost:8080/paragliding/api/", code: 200},
		{method: "GET", url: "http://localhost:8080/paragliding/api/rr", code: 404},
		{method: "POST", url: "http://localhost:8080/paragliding/api/", code: 405},
	}
	for _, testCase := range TestTable {
		req, err := http.NewRequest(testCase.method, testCase.url, nil)
		if err != nil {
			t.Errorf("Unexpected error, %d", err)
		}

		writer := httptest.NewRecorder()
		Index(writer, req)
		if writer.Code != testCase.code {
			t.Errorf("got resp code %v, expected %v", writer.Code, testCase.code)
		}
	}
}

//Testing the func RegAndShowTrack
func TestRegAndShowTrack(t *testing.T) {
	Connect()
	TestTable := []struct {
		method string
		url    string
		code   int
	}{
		{method: "GET", url: "localhost:8080/paragliding/api/track", code: 200},
		{method: "GET", url: "localhost:8080/paragliding/api/track", code: 200},
		{method: "PUT", url: "localhost:8080/paragliding/api/track123", code: 405},
	}
	for _, testCase := range TestTable {
		req, err := http.NewRequest(testCase.method, testCase.url, nil)
		if err != nil {
			t.Errorf("Unexpected error, %d", err)
		}

		writer := httptest.NewRecorder()
		RegAndShowTrack(writer, req)
		if writer.Code != testCase.code {
			t.Errorf("got resp code %v, expected %v", writer.Code, testCase.code)
		}
		if testCase.method == "GET" {
			var ids []int
			if err := json.NewDecoder(writer.Body).Decode(&ids); err != nil {
				t.Errorf("Error reading the body %d", err)
			}
			if len(ids) < 1 {
				t.Errorf("Expected atleast one track to be registered, got %v", len(ids))
			}
		}
	}
}

//Testing the function ShowTrackInfo
func TestShowTrackInfo(t *testing.T) {
	Connect()
	TestTable := []struct {
		method string
		url    string
		code   int
	}{
		{method: "GET", url: "localhost:8080/paragliding/api/track/1", code: 404},
		{method: "GET", url: "localhost:8080/paragliding/api/track/2", code: 404},
		{method: "PUT", url: "localhost:8080/paragliding/api/track/1", code: 404},
	}
	for _, testCase := range TestTable {
		req, err := http.NewRequest(testCase.method, testCase.url, nil)
		if err != nil {
			t.Errorf("Unexpected error, %d", err)
		}

		writer := httptest.NewRecorder()
		ShowTrackInfo(writer, req)
		if writer.Code != testCase.code {
			t.Errorf("got resp code %v, expected %v", writer.Code, testCase.code)
		}
		if testCase.code == 200 {
			track := database.Track{}
			if err := json.NewDecoder(writer.Body).Decode(&track); err != nil {
				t.Errorf("Error reading the body %v", err)
			}
		}
	}
}

// Testing the function GetTickerInfo
func TestGetTickerInfo(t *testing.T) {
	Connect()
	TestTable := []struct {
		method string
		url    string
		code   int
	}{
		{method: "GET", url: "localhost:8080/api/ticker/?cap=4", code: 200},
		{method: "GET", url: "localhost:8080/api/ticker/1231/asda?cap=4", code: 404},
		{method: "POST", url: "localhost:8080/api/ticker/?cap=4", code: 405},
		{method: "PUT", url: "localhost:8080/api/ticker/?cap=4", code: 405},
	}
	for _, testCase := range TestTable {
		req, err := http.NewRequest(testCase.method, testCase.url, nil)
		if err != nil {
			t.Errorf("Unexpected error, %d", err)
		}

		writer := httptest.NewRecorder()
		GetTickerInfo(writer, req)
		if writer.Code != testCase.code {
			t.Errorf("got resp code %v, expected %v", writer.Code, testCase.code)
		}
		if writer.Code == 200 {
			ticker := database.Ticker{}
			if err := json.NewDecoder(writer.Body).Decode(&ticker); err != nil {
				t.Errorf("Error reading the body %v", err)
			}
		}
	}
}

// Testing GetLatestTicker
func TestGetLatestTicker(t *testing.T) {
	Connect()
	TestTable := []struct {
		method string
		url    string
		code   int
	}{
		{method: "GET", url: "localhost:8080/api/ticker/latest", code: 200},
		{method: "GET", url: "localhost:8080/api/ticker/latest/", code: 404},
		{method: "POST", url: "localhost:8080/api/ticker/latest", code: 405},
		{method: "PUT", url: "localhost:8080/api/ticker/latest", code: 405},
	}
	for _, testCase := range TestTable {
		req, err := http.NewRequest(testCase.method, testCase.url, nil)
		if err != nil {
			t.Errorf("Unexpected error, %d", err)
		}

		writer := httptest.NewRecorder()
		GetLatestTicker(writer, req)
		if writer.Code != testCase.code {
			t.Errorf("got resp code %v, expected %v", writer.Code, testCase.code)
		}
		if writer.Code == 200 {
			//var Last time.Time
			last := writer.Body.String()
			if last == "" {
				t.Errorf("Error reading the body %v", last)
			}
		}
	}
}

// Testing the func PopulateTickerInfo

func TestPopulateTickerInfo(t *testing.T) {
	Connect()
	err := PopulateTickerInfo(5, false, "")
	if err != nil {
		t.Errorf("Did not expect error, got %v", err)
	}

	err = PopulateTickerInfo(5, true, "2018-10-26T17:12:37.128Z")
	if err != nil {
		t.Errorf("Did not expect error, got %v", err)
	}

	//Expecting it to fail
	err = PopulateTickerInfo(5, true, "2018-10-26T17:12:37.128Z09000")
	if err == nil {
		t.Errorf("was expecting error, got %v", err)
	}
}

// Testing the admin handler for tracks_count (the func GetTracksCount)
func TestGetTracksCount(t *testing.T) {
	Connect()
	TestTable := []struct {
		method string
		url    string
		code   int
	}{
		{method: "GET", url: "localhost:8080/admin/api/tracks_count?admincode=12345", code: 200},
		{method: "GET", url: "localhost:8080/admin/api/tracks_count?admincode=1234", code: 404},
		{method: "POST", url: "localhost:8080/admin/api/tracks_count?admincode=12345", code: 405},
		{method: "PUT", url: "localhost:8080/admin/api/tracks_count?admincode=12345", code: 405},
	}
	for _, testCase := range TestTable {
		req, err := http.NewRequest(testCase.method, testCase.url, nil)
		if err != nil {
			t.Errorf("Unexpected error, %d", err)
		}

		writer := httptest.NewRecorder()
		GetTracksCount(writer, req)
		if writer.Code != testCase.code {
			t.Errorf("got resp code %v, expected %v", writer.Code, testCase.code)
		}
		if writer.Code == 200 {
			//var Last time.Time
			last := writer.Body.String()
			if last == "" {
				t.Errorf("Error reading the body %v", last)
			}
		}
	}
}

// Testing the admin handler for deleting all the tracks (func DeleteAllTracks)

func TestDeleteAllTracks(t *testing.T) {
	Connect()
	TestTable := []struct {
		method string
		url    string
		code   int
	}{
		{method: "DELETE", url: "localhost:8080/admin/api/tracks?admincode=12345", code: 200},
		{method: "DELETE", url: "localhost:8080/admin/api/tracks?admincode=1234", code: 404},
		{method: "DELETE", url: "localhost:8080/admin/api/tracks", code: 404},
		{method: "POST", url: "localhost:8080/admin/api/tracks?admincode=12345", code: 405},
		{method: "PUT", url: "localhost:8080/admin/api/tracks?admincode=12345", code: 405},
	}
	for _, testCase := range TestTable {
		req, err := http.NewRequest(testCase.method, testCase.url, nil)
		if err != nil {
			t.Errorf("Unexpected error, %d", err)
		}

		writer := httptest.NewRecorder()
		DeleteAllTracks(writer, req)
		if writer.Code != testCase.code {
			t.Errorf("got resp code %v, expected %v", writer.Code, testCase.code)
		}
		if writer.Code == 200 {
			last := writer.Body.String()
			if last == "" {
				t.Errorf("Error reading the body %v", last)
			}
			if !(last == "OK!" || last == "Empty record") {
				t.Errorf("Expecting OK!, got %v", last)
			}
		}
	}
}

// Testing the webhook endpoint RegisterWebhook
func TestRegisterWebhook(t *testing.T) {
	Connect()
	TestTable := []struct {
		method string
		url    string
		code   int
	}{
		{method: "GET", url: "localhost:8080/api/webhook/new_track/", code: 404},
		{method: "GET", url: "localhost:8080/api/webhook/new_track/100", code: 404},
		{method: "PUT", url: "localhost:8080/api/webhook/new_track/", code: 405},
	}
	for _, testCase := range TestTable {
		req, err := http.NewRequest(testCase.method, testCase.url, nil)
		if err != nil {
			t.Errorf("Unexpected error, %d", err)
		}

		writer := httptest.NewRecorder()
		RegisterWebhook(writer, req)
		if writer.Code != testCase.code {
			t.Errorf("got resp code %v, expected %v", writer.Code, testCase.code)
		}
		if writer.Code == 200 {

		}
	}
}
