package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Meta ...
type Meta struct {
	Uptime  string `json:"uptime"`
	Info    string `json:"info"`
	Version string `json:"version"`
}

// apiRedirect redirects the user to the api
func apiRedirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/paragliding/api", http.StatusSeeOther)
}

// apiGet gets the meta information about the api
func apiGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Using a function from stackoverflow I get the difference between two dates and format them according to ISO8601
	y, mo, d, h, m, s := diff(startTime, time.Now())

	// Creates a string with the uptime in ISO8601 format
	uptime := fmt.Sprintf("P%vY%vM%vDT%vH%vM%vS", y, mo, d, h, m, s)

	// Creates a new meta struct to get the data in json format
	meta := Meta{uptime, "Service for Paragliding tracks.", "v1"}
	json.NewEncoder(w).Encode(meta)
}
