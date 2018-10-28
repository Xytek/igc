package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// Ticker is used to sort through tracks from the database
type Ticker struct {
	Latest     int64         `json:"t_latest"`
	Start      int64         `json:"t_start"`
	Stop       int64         `json:"t_stop"`
	Tracks     []int         `json:"tracks"`
	Processing time.Duration `json:"processing"`
}

const pagingCap = 5

// tickerLatestGet returns the timestamp of the latest added track
func tickerLatestGet(w http.ResponseWriter, r *http.Request) {
	// Brings out the latest Track from the db
	t, ok := db.getTrack(IDs[len(IDs)-1])
	if !ok {
		http.Error(w, "There was an error when trying to Get() the id", http.StatusNotFound)
	}

	json.NewEncoder(w).Encode(t.Timestamp)
}

// tickerLatestGet returns the timestamp of the latest added track
func tickerGet(w http.ResponseWriter, r *http.Request) {
	// Tracks the start time of the function for processing information
	funcStart := time.Now()

	var ticker Ticker
	var tracks = db.getAllTracks()

	// Keeps the amount of tracks
	var trackCount = len(tracks)
	if trackCount < 0 {
		http.Error(w, "No tracks have been added yet", http.StatusBadRequest)
		return
	}

	// Checks if there are enough tracks to fill the pagingCap, and sets the limit lower otherwise
	var limit int
	if trackCount >= pagingCap {
		limit = pagingCap
	} else {
		limit = trackCount
	}

	// Fill start/latest with the first/last timestamp respectively
	ticker.Start = tracks[0].Timestamp
	ticker.Latest = tracks[trackCount-1].Timestamp

	// Fills the pager with an appropiate amount of track IDs
	for i := 0; i < limit; i++ {
		ticker.Tracks = append(ticker.Tracks, tracks[i].SimpleID)
	}

	// Find the last ID of the paged ones and get the timestamp of the track with that id
	t, _ := db.getTrack(ticker.Tracks[len(ticker.Tracks)-1])
	ticker.Stop = t.Timestamp

	// Checks how long it's been since the start of the function in ms
	ticker.Processing = time.Since(funcStart) / 1000000

	// Returns all information about the ticker
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ticker)
}

func tickerTimestampGet(w http.ResponseWriter, r *http.Request) {
	// Get ID from the URL
	ts := mux.Vars(r)["timestamp"]
	i, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		http.Error(w, "Not a valid timestamp", http.StatusBadRequest)
		return
	}

	// Tracks the start time of the function for processing information
	funcStart := time.Now()

	var ticker Ticker
	tracks := db.getTracksAfter(i)

	// Keeps the amount of tracks
	var trackCount = len(tracks)
	if trackCount < 0 {
		http.Error(w, "No tracks have been added since this timestamp", http.StatusBadRequest)
		return
	}

	// Checks if there are enough tracks to fill the pagingCap, and sets the limit lower otherwise
	var limit int
	if trackCount >= pagingCap {
		limit = pagingCap
	} else {
		limit = trackCount
	}

	// Fill start/latest with the first/last timestamp respectively
	ticker.Start = tracks[0].Timestamp
	ticker.Latest = tracks[trackCount-1].Timestamp

	// Fills the pager with an appropiate amount of track IDs
	for i := 0; i < limit; i++ {
		ticker.Tracks = append(ticker.Tracks, tracks[i].SimpleID)
	}

	// Find the last ID of the paged ones and get the timestamp of the track with that id
	t, _ := db.getTrack(ticker.Tracks[len(ticker.Tracks)-1])
	ticker.Stop = t.Timestamp

	// Checks how long it's been since the start of the function in ms
	ticker.Processing = time.Since(funcStart) / 1000000

	// Returns all information about the ticker
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ticker)
}
