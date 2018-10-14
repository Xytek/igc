package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	igc "github.com/marni/goigc"
)

/*
Test overview:
1. Not implemented commands
2. Malformed URLs
3. Get all track IDs when empty
4. Get all track IDs when filled
5. Get specific track by ID
6. Get specific pilot by ID
7. POST with both good and bad URL
8. POST checking response and adding of new track and id
9. Get meta data
10. checkURL with good and bad values
*/

func Test_handlerAPIIGC_notImplemented(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handlerAPIIGC))
	defer ts.Close()

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodDelete, ts.URL, nil)
	if err != nil {
		t.Errorf("Error constructing the DELETE request %s", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Error making the DELETE request %s", err)
	}
	if resp.StatusCode != http.StatusNotImplemented {
		t.Errorf("Expected StatusCode %d, received %d", http.StatusNotImplemented, resp.StatusCode)
	}
}

func Test_handlerAPIIGCMORE_malformedURL(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handlerAPIIGCMORE))
	defer ts.Close()
	testCases := []string{
		ts.URL,
		ts.URL + "/",
		ts.URL + "/igcinfo",
		ts.URL + "/igcinfo/",
		ts.URL + "/igcinfo/api/",
		ts.URL + "/igcinfo/api/igc/",
		ts.URL + "/igcinfo/api/igc/ID",
		ts.URL + "/igcinfo/api/igc/ID/",
		ts.URL + "/igcinfo/api/igc/ID/piloti",
		ts.URL + "/igcinfo/api/igc/ID/pilot/",
	}
	for _, tstring := range testCases {
		resp, err := http.Get(tstring)
		if err != nil {
			t.Errorf("Error making the GET request, %s", err)
		}
		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("For route. %s, expected StatusCode %d, received %d\n", tstring, http.StatusNotFound, resp.StatusCode)
			return
		}
	}
}

func Test_handlerAPIIGC_getAllTracks_empty(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handlerAPIIGC))
	defer ts.Close()
	resp, err := http.Get(ts.URL + "/igcinfo/api/igc")
	if err != nil {
		t.Errorf("Error making the GET request, %s", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected StatusCode %d, received %d", http.StatusOK, resp.StatusCode)
		return
	}
	var a []interface{}
	err = json.NewDecoder(resp.Body).Decode(&a)
	if err != nil {
		t.Errorf("Error parsing the expected JSON body. Got %s", err)
	}
	if len(a) != 0 {
		t.Errorf("Expected empty array, got %s", a)
	}
}

func Test_handlerAPIIGCMORE_getID_ID0(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handlerAPIIGCMORE))
	defer ts.Close()
	db = TrackDB{}
	db.Init()
	track, _ := igc.ParseLocation("http://skypolaris.org/wp-content/uploads/IGS%20Files/Madrid%20to%20Jerez.igc")
	totalDistance := 0.0
	for i := 0; i < len(track.Points)-1; i++ {
		totalDistance += track.Points[i].Distance(track.Points[i+1])
	}
	testTrack := Track{track.Header.Date, track.Pilot, track.GliderType, track.GliderID, totalDistance}
	db.Add(testTrack, ID{"ID0"})

	resp, err := http.Get(ts.URL + "/igcinfo/api/igc/id0")
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
		t.Errorf("Error parsing the expected JSON body. Got %s", err)
	}
	if a.HDate != testTrack.HDate || a.Pilot != testTrack.Pilot || a.Glider != testTrack.Glider || a.GliderID != testTrack.GliderID || a.TrackLength != testTrack.TrackLength {
		t.Errorf("Tracks do not match! Got: %v, expected %v\n", a, testTrack)
	}
}

func Test_handlerAPIIGCMORE_getPilot(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handlerAPIIGCMORE))
	defer ts.Close()
	db = TrackDB{}
	db.Init()
	track, _ := igc.ParseLocation("http://skypolaris.org/wp-content/uploads/IGS%20Files/Madrid%20to%20Jerez.igc")
	totalDistance := 0.0
	for i := 0; i < len(track.Points)-1; i++ {
		totalDistance += track.Points[i].Distance(track.Points[i+1])
	}
	testTrack := Track{track.Header.Date, track.Pilot, track.GliderType, track.GliderID, totalDistance}
	db.Add(testTrack, ID{"ID0"})

	resp, err := http.Get(ts.URL + "/igcinfo/api/igc/id0/pilot")
	if err != nil {
		t.Errorf("Error making the GET request, %s", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected StatusCode %d, received %d", http.StatusOK, resp.StatusCode)
		return
	}
	a, _ := ioutil.ReadAll(resp.Body)
	if string(a) != testTrack.Pilot {
		t.Errorf("Expected %s got %s", track.Pilot, a)
	}
}
func Test_handlerAPIIG_getID_allIDs(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handlerAPIIGC))
	defer ts.Close()
	db = TrackDB{}
	db.Init()
	IDs = nil // Make sure there's no leftovers from previous tests
	db.Add(Track{}, ID{"ID0"})
	db.Add(Track{}, ID{"ID1"})

	resp, err := http.Get(ts.URL + "/igcinfo/api/igc")
	if err != nil {
		t.Errorf("Error making the GET request, %s", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected StatusCode %d, received %d", http.StatusOK, resp.StatusCode)
		return
	}
	var a []string
	err = json.NewDecoder(resp.Body).Decode(&a)
	if err != nil {
		t.Errorf("Error parsing the expected JSON body. Got %s", err)
	}
	if len(a) != 2 {
		t.Errorf("Expected array with two elements, got %v", a)
	}
	if a[0] != "ID0" || a[1] != "ID1" {
		t.Errorf("The first ID is supposed to be ID0 but is %v. The second should be ID1 but is %v", a[0], a[1])
	}
}

func Test_handler_POST(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handlerAPIIGC))
	defer ts.Close()

	db = TrackDB{}
	db.Init()

	// Test empty body
	resp, err := http.Post(ts.URL+"/igcinfo/api/igc", "application/json", nil)
	if err != nil {
		t.Errorf("Error creating the POST request, %s", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected StatusCode %d, received %d", http.StatusBadRequest, resp.StatusCode)
	}

	// Test proper body
	test := "{\"url\":\"http://skypolaris.org/wp-content/uploads/IGS%20Files/Madrid%20to%20Jerez.igc\"}"
	resp, err = http.Post(ts.URL+"/igcinfo/api/igc", "application/json", strings.NewReader(test))
	if err != nil {
		t.Errorf("Error creating the POST request, %s", err)
	}
	if resp.StatusCode != http.StatusOK {
		all, _ := ioutil.ReadAll(resp.Body)
		t.Errorf("Expected StatusCode %d, received %d, Body: %s", http.StatusOK, resp.StatusCode, all)
	}

	// Test malformed URL body
	badTest := "{\"url\":\"http://skypolaris.org/wp-content/uploads/IGS%20Files/Maadrid%20to%20Jerez.igc\"}"

	resp, err = http.Post(ts.URL+"/igcinfo/api/igc", "application/json", strings.NewReader(badTest))
	if err != nil {
		t.Errorf("Error creating the POST request, %s", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		all, _ := ioutil.ReadAll(resp.Body)
		t.Errorf("Expected StatusCode %d, received %d, Body: %s", http.StatusBadRequest, resp.StatusCode, all)
	}
}

func Test_handler_POST_ResponseAndGet(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handlerAPIIGC))
	defer ts.Close()
	db = TrackDB{}
	db.Init()

	// Test proper body
	test := "{\"url\":\"http://skypolaris.org/wp-content/uploads/IGS%20Files/Madrid%20to%20Jerez.igc\"}"
	track, _ := igc.ParseLocation("http://skypolaris.org/wp-content/uploads/IGS%20Files/Madrid%20to%20Jerez.igc")
	resp, _ := http.Post(ts.URL+"/igcinfo/api/igc", "application/json", strings.NewReader(test))
	var a ID
	err := json.NewDecoder(resp.Body).Decode(&a)
	if err != nil {
		t.Errorf("Error parsing the expected JSON body. Got %s", err)
	}
	tr, erro := db.Get(a.ID)
	if !erro {
		t.Errorf("The track was not added properly, it could not be found on %s", a.ID)
	}
	if tr.HDate != track.Date || tr.Pilot != track.Pilot || tr.Glider != track.GliderType || tr.GliderID != track.GliderID {
		t.Error("The track was not added properly. Its data is incorrect")
	}
}

func Test_handlerAPI_GETMeta(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handlerAPI))
	defer ts.Close()
	y, mo, d, h, m, s := diff(startTime, time.Now())
	uptime := fmt.Sprintf("P%vY%vM%vDT%vH%vM%vS", y, mo, d, h, m, s)
	meta := Meta{uptime, "Service for IGC tracks.", "v1"}
	resp, err := http.Get(ts.URL + "/igcinfo/api")
	if err != nil {
		t.Errorf("Error making the GET request, %s", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected StatusCode %d, received %d", http.StatusOK, resp.StatusCode)
		return
	}
	var a Meta
	err = json.NewDecoder(resp.Body).Decode(&a)
	if err != nil {
		t.Errorf("Error parsing the expected JSON body. Got %s", err)
	}
	if a.Version != meta.Version || a.Info != meta.Info || a.Uptime != meta.Uptime {
		t.Errorf("Expected %v got %v", meta, a)
	}
}

func Test_checkURL(t *testing.T) {
	goodCases := []string{
		"http://skypolaris.org/wp-content/uploads/IGS%20Files/Madrid%20to%20Jerez.igc",
		"http://skypolaris.org/wp-content/uploads/IGS%20Files/Jarez%20to%20Senegal.igc",
		"http://skypolaris.org/wp-content/uploads/IGS%20Files/Boavista%20Medellin.igc",
		"http://skypolaris.org/wp-content/uploads/IGS%20Files/Medellin%20Guatemala.igc",
	}
	for _, tstring := range goodCases {
		if checkURL(tstring) == false {
			t.Errorf("A legit URL could not be read: %s", tstring)
		}
	}
	badCases := []string{
		"http://skypolaris.org/wp-content/uploads/IGS%20Files/Madrid%20to%20Jerez",
		"http://skypolaris.org/wp-content/uploads/IGS%20Files/Jarez%2To%2Senegal.igc",
		"http://skypolaris.com/wp-content/uploads/IGS%20Files/Boavista%20Medellin.igc",
		"skypolaris.org/wp-content/uploads/IGS%20Files/Medellin%20Guatemala.igc",
	}
	for _, tstring := range badCases {
		if checkURL(tstring) != false {
			t.Errorf("An illegit URL could be read: %s", tstring)
		}
	}
}
