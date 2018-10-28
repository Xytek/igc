package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	igc "github.com/marni/goigc"
)

func Test_trackGet(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(trackGet))
	defer ts.Close()
	DBInitTest()

	tData := Track{"A date", "A pilot", "type", "GliderID", 0, "url", 0, 0}
	db.addTrack(tData)
	db.addTrack(tData)

	resp, err := http.Get(ts.URL + "/paragliding/api/track")
	if err != nil {
		t.Errorf("Error making the GET request, %s", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected StatusCode %d, received %d", http.StatusOK, resp.StatusCode)
		return
	}
	var a []int
	err = json.NewDecoder(resp.Body).Decode(&a)
	if len(a) != 2 {
		t.Errorf("Expected two elements, got %v", len(a))
	}
	db.deleteAllTracks()
}

func Test_trackPost(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(trackPost))
	defer ts.Close()
	DBInitTest()

	// Test empty body
	resp, err := http.Post(ts.URL+"/paragliding/api/track", "application/json", nil)
	if err != nil {
		t.Errorf("Error creating the POST request, %s", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected StatusCode %d, received %d", http.StatusBadRequest, resp.StatusCode)
	}

	// Test malformed URL body
	badTest := "{\"url\":\"http://skypolaris.org/wp-content/uploads/IGS%20Files/Madrid%20to%20Jerez.ic\"}"

	resp, err = http.Post(ts.URL+"/paragliding/api/track", "application/json", strings.NewReader(badTest))
	if err != nil {
		t.Errorf("Error creating the POST request, %s", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		all, _ := ioutil.ReadAll(resp.Body)
		t.Errorf("Expected StatusCode %d, received %d, Body: %s", http.StatusBadRequest, resp.StatusCode, all)
	}

	// Test proper body, actual data input and response
	track, _ := igc.ParseLocation("http://skypolaris.org/wp-content/uploads/IGS%20Files/Madrid%20to%20Jerez.igc")
	test := "{\"url\":\"http://skypolaris.org/wp-content/uploads/IGS%20Files/Madrid%20to%20Jerez.igc\"}"
	resp, err = http.Post(ts.URL+"/paragliding/api/track", "application/json", strings.NewReader(test))
	if err != nil {
		t.Errorf("Error creating the POST request, %s", err)
	}
	if resp.StatusCode != http.StatusOK {
		all, _ := ioutil.ReadAll(resp.Body)
		t.Errorf("Expected StatusCode %d, received %d, Body: %s", http.StatusOK, resp.StatusCode, all)
	}
	var a ID
	err = json.NewDecoder(resp.Body).Decode(&a)
	if err != nil {
		t.Errorf("Error parsing the expected JSON body. Got %s", err)
	}
	tr, erro := db.getTrack(a.ID)
	if !erro {
		t.Errorf("The track was not added properly, it could not be found on %v", a.ID)
	}
	if tr.HDate != track.Date.String() || tr.Pilot != track.Pilot || tr.Glider != track.GliderType || tr.GliderID != track.GliderID {
		t.Error("The track was not added properly. Its data is incorrect")
	}
	db.deleteAllTracks()
}

func Test_trackIDGet(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/paragliding/api/track/{igcId}", trackIDGet).Methods("GET")
	ts := httptest.NewServer(r)
	defer ts.Close()
	DBInitTest()

	tData := Track{"A date", "Pilot", "type", "GliderID", 0, "url", 0, 0}
	db.addTrack(tData)

	resp, err := http.Get(ts.URL + "/paragliding/api/track/0")
	if err != nil {
		t.Errorf("Error making the GET request, %s", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected StatusCode %d, received %d", http.StatusOK, resp.StatusCode)
		return
	}
	var a Track
	err = json.NewDecoder(resp.Body).Decode(&a)
	if err != nil {
		t.Errorf("Problem decoding: %v", err)
	}
	if a.Pilot != tData.Pilot {
		t.Errorf("Expected %s got %s", tData.Pilot, a.Pilot)
	}
	db.deleteAllTracks()
}

func Test_trackIDFieldGet(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/paragliding/api/track/{igcId}/{igcField}", trackIDFieldGet).Methods("GET")
	ts := httptest.NewServer(r)
	defer ts.Close()
	DBInitTest()

	tData := Track{"A date", "Pilot", "type", "GliderID", 0, "url", 0, 0}
	db.addTrack(tData)

	resp, err := http.Get(ts.URL + "/paragliding/api/track/0/pilot")
	if err != nil {
		t.Errorf("Error making the GET request, %s", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected StatusCode %d, received %d", http.StatusOK, resp.StatusCode)
		return
	}
	a, _ := ioutil.ReadAll(resp.Body)
	if string(a) != tData.Pilot {
		t.Errorf("Expected %s got %s", tData.Pilot, a)
	}
	db.deleteAllTracks()
}
