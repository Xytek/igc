package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	igc "github.com/marni/goigc"
)

// Track holds all the relevant information about a track in the database
type Track struct {
	HDate       string  `json:"H_date"`
	Pilot       string  `json:"pilot"`
	Glider      string  `json:"glider"`
	GliderID    string  `json:"glider_id"`
	TrackLength float64 `json:"track_length"`
	URL         string  `json:"track_src_url"`
	Timestamp   int64   `json:"-" bson:"timestamp"`
	SimpleID    int     `json:"-" bson:"simpleid"`
}

// ID is used to return json after POSTing URL
type ID struct {
	ID int `json:"id"`
}

// URL is used to read in an url in json format
type URL struct {
	URL string `json:"url"`
}

// trackPost reads in an URL, gets the IGC information from it and adds it to the database
func trackPost(w http.ResponseWriter, r *http.Request) {
	//Decode incoming url
	var u URL
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// checkURL uses regex to check if correct format
	if checkURL(u.URL) == false {
		http.Error(w, "The URL is formated incorrectly", http.StatusBadRequest)
		return
	}

	// Gets api information from the URL
	track, err := igc.ParseLocation(u.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Calculates the total distance between the first and last point, going through each point
	totalDistance := 0.0
	for i := 0; i < len(track.Points)-1; i++ {
		totalDistance += track.Points[i].Distance(track.Points[i+1])
	}

	// Creates a new Track with data from the API
	var t Track
	t.HDate = track.Date.String()
	t.Pilot = track.Pilot
	t.Glider = track.GliderType
	t.GliderID = track.GliderID
	t.TrackLength = totalDistance
	t.URL = u.URL

	// An ID struct in order to return ID as JSON
	var i ID
	i.ID = lastUsed

	//Pass DB Credentials and new track to function
	t = db.addTrack(t)

	// Write back the ID
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(i)
}

// trackGet returns a json array with all internal IDs in it
func trackGet(w http.ResponseWriter, r *http.Request) {
	// To avoid a null value and ensure we get an empty array instead
	if len(IDs) == 0 {
		IDs = make([]int, 0)
	}

	json.NewEncoder(w).Encode(IDs)
}

// trackIDGet finds and returnes a json array with information about one track whos ID is found in the URL
func trackIDGet(w http.ResponseWriter, r *http.Request) {
	// Get ID from the URL
	id := mux.Vars(r)["igcId"]

	// Convert id to integer
	i, err := strconv.Atoi(id)
	if err != nil || i < 0 {
		http.Error(w, "Malformed url. The ID must be a positive integer", http.StatusNotFound)
	}

	// Loops through the array of all IDs to check if the one in the url is among them
	idExists := false
	for j := range IDs {
		if j == i {
			idExists = true
			break
		}
	}

	// If the ID in the URL does not exist for whatever reason then this is thrown
	if !idExists {
		http.Error(w, "Malformed url. The ID does not exist.", http.StatusNotFound)
		return
	}

	// Brings out the Track with said ID from the db
	t, ok := db.getTrack(i)
	if !ok {
		http.Error(w, "There was an error when trying to Get() the id", http.StatusNotFound)
		return
	}

	// Writes out the track with the correct ID
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}

// trackIDFieldGet finds and returns information about a single attribute of a track
func trackIDFieldGet(w http.ResponseWriter, r *http.Request) {
	// Get ID from the URL
	vari := mux.Vars(r)
	id := vari["igcId"]
	field := vari["igcField"]

	// Convert id to integer
	i, err := strconv.Atoi(id)
	if err != nil || i < 0 {
		http.Error(w, "Malformed url. The ID must be a positive integer", http.StatusNotFound)
		return
	}

	// Loops through the array of all IDs to check if the one in the url is among them
	idExists := false
	for j := range IDs {
		if j == i {
			idExists = true
			break
		}
	}

	// If the ID in the URL does not exist for whatever reason then this is thrown
	if !idExists {
		http.Error(w, "Malformed url. The ID does not exist.", http.StatusNotFound)
		return
	}

	// Brings out the Track with said ID from the db
	t, ok := db.getTrack(i)
	if !ok {
		http.Error(w, "There was an error when trying to Get() the id", http.StatusNotFound)
	}

	// Switch that goes through the fields that we can look up
	switch strings.ToUpper(field) {
	case "PILOT":
		fmt.Fprint(w, t.Pilot)
	case "GLIDER":
		fmt.Fprint(w, t.Glider)
	case "GLIDER_ID":
		fmt.Fprint(w, t.GliderID)
	case "TRACK_LENGTH":
		fmt.Fprint(w, t.TrackLength)
	case "H_DATE":
		fmt.Fprint(w, t.HDate)
	case "TRACK_SRC_URL":
		fmt.Fprint(w, t.URL)
	default:
		http.Error(w, "Malformed URL. Your options are pilot, glider, glider_id, track_length, h_date or track_src_url", http.StatusNotFound)
		return
	}
}
