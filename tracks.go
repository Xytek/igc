package main

import (
	"time"
)

// IDs is an array holding all existing IDs as strings
var IDs []string

// Meta ...
type Meta struct {
	Uptime  string `json:"uptime"`
	Info    string `json:"info"`
	Version string `json:"version"`
}

// URL ...
type URL struct {
	URL string `json:"url"`
}

// ID ...
type ID struct {
	ID string `json:"id"`
}

// Track ...
type Track struct {
	HDate       time.Time `json:"H_date"`
	Pilot       string    `json:"pilot"`
	Glider      string    `json:"glider"`
	GliderID    string    `json:"glider_id"`
	TrackLength float64   `json:"track_length"`
}

// TrackDB ....
type TrackDB struct {
	tracks map[string]Track
}

// Init ...
func (db *TrackDB) Init() {
	db.tracks = make(map[string]Track)
}

// Add ...
func (db *TrackDB) Add(t Track, i ID) {
	db.tracks[i.ID] = t
	IDs = append(IDs, i.ID)
}

// Count ...
func (db *TrackDB) Count() int {
	return len(db.tracks)
}

// Get ....
func (db *TrackDB) Get(keyID string) (Track, bool) {
	t, ok := db.tracks[keyID]
	return t, ok
}
