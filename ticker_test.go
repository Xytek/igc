package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func Test_tickerGet(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(tickerGet))
	defer ts.Close()
	DBInitTest()

	tData := Track{"A date", "A pilot", "type", "GliderID", 0, "url", 0, 0}
	db.addTrack(tData)
	db.addTrack(tData)

	resp, err := http.Get(ts.URL + "/paragliding/api/ticker")
	if err != nil {
		t.Errorf("Error making the GET request, %s", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected StatusCode %d, received %d", http.StatusOK, resp.StatusCode)
		return
	}
	var a Ticker
	err = json.NewDecoder(resp.Body).Decode(&a)
	if len(a.Tracks) != 2 {
		t.Errorf("Expected two elements, got %v", len(a.Tracks))
	}
	db.deleteAllTracks()
}

func Test_tickerLatestGet(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(tickerLatestGet))
	defer ts.Close()
	DBInitTest()

	tData := Track{"A date", "A pilot", "type", "GliderID", 0, "url", 0, 0}
	tData = db.addTrack(tData)

	resp, err := http.Get(ts.URL + "/paragliding/api/ticker/latest")
	if err != nil {
		t.Errorf("Error making the GET request, %s", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected StatusCode %d, received %d", http.StatusOK, resp.StatusCode)
		return
	}
	var a int64
	err = json.NewDecoder(resp.Body).Decode(&a)
	if a != tData.Timestamp {
		t.Errorf("Expected %v, got %v", tData.Timestamp, a)
	}
	db.deleteAllTracks()
}

func Test_tickerTimestampGet(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/paragliding/api/ticker/{timestamp}", tickerTimestampGet).Methods("GET")
	ts := httptest.NewServer(r)
	defer ts.Close()
	DBInitTest()

	tData := Track{"A date", "A pilot", "type", "GliderID", 0, "url", 0, 0}
	db.addTrack(tData)
	db.addTrack(tData)

	resp, err := http.Get(ts.URL + "/paragliding/api/ticker/50000")
	if err != nil {
		t.Errorf("Error making the GET request, %s", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected StatusCode %d, received %d", http.StatusOK, resp.StatusCode)
		return
	}
	var a Ticker
	err = json.NewDecoder(resp.Body).Decode(&a)
	if len(a.Tracks) != 2 {
		t.Errorf("Expected two elements, got %v", len(a.Tracks))
	}
	db.deleteAllTracks()
}
