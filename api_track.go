package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/marni/goigc"
)

var db TrackDB

// handlerAPI writes out meta data
func handlerAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Using a function from stackoverflow I get the difference between two dates and format them according to ISO8601
	y, mo, d, h, m, s := diff(startTime, time.Now())
	uptime := fmt.Sprintf("P%vY%vM%vDT%vH%vM%vS", y, mo, d, h, m, s)
	meta := Meta{uptime, "Service for IGC tracks.", "v1"} // Creates a new meta struct to get the data in json format
	json.NewEncoder(w).Encode(meta)
}

// handlerAPIIGC deals with POSTing URL, creating new Tracks/IDs and giving appropiate response on both POST and GET
func handlerAPIIGC(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case "POST":
		if r.Body == nil {
			http.Error(w, "Track POST request must have JSON body", http.StatusBadRequest)
			return
		}
		// Reads in the URL
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
		var i ID
		i.ID = ("ID" + strconv.Itoa(lastUsed))
		t := Track{track.Header.Date,
			track.Pilot,
			track.GliderType,
			track.GliderID,
			totalDistance}
		lastUsed++ // Increments the global variable to keep IDs unique

		// Initiates the db if it hasn't been already
		if db.tracks == nil {
			db.Init()
		}
		// Adds the newly created Track and ID to the db
		db.Add(t, i)
		// Obligatory response
		json.NewEncoder(w).Encode(i)
		return
	case "GET":
		// To avoid a null value and ensure we get an empty array instead
		if len(IDs) == 0 {
			IDs = make([]string, 0)
		}
		json.NewEncoder(w).Encode(IDs)
		return
	default:
		http.Error(w, "Not implemented yet", http.StatusNotImplemented)
		return
	}
}

func handlerAPIIGCMORE(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	idExists := false
	// Loops through the array of all IDs to check if the one in the url is among them
	for i := 0; i < len(IDs); i++ {
		if IDs[i] == strings.ToUpper(parts[4]) {
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
	t, ok := db.Get(strings.ToUpper(parts[4]))
	// Just a last check, but this really should never happen at this point
	if !ok {
		http.Error(w, "This really shouldn't have happened", http.StatusNotFound)
	}
	// It's already checked if the id is good, so now it makes sure there's an appropiate amount of parts
	if len(parts) == 5 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(t)
	}

	// Here's all the logic for the <id>/<field> part.
	// It has a case for each of the possibilities from the task, anything else will go 404
	if len(parts) == 6 {
		switch strings.ToUpper(parts[5]) {
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
		default:
			http.Error(w, "Malformed URL. Your options are pilot, glider, glider_id, track_length and H_date", http.StatusNotFound)
			return
		}

	}
	// Bonus message in case the user tries any more / after <field>
	if len(parts) > 6 {
		http.Error(w, "Malformed url. You have one or more / too much.", http.StatusNotFound)
	}
}

// checkURL uses a regular expression to check that the URL is of the type that this program can read.
func checkURL(u string) bool {
	// There are two types of files that can be read. Only difference is one uses %20to%20 and other uses %20 between the cities. Therefore the rest is hardcoded
	check, _ := regexp.MatchString("^(http://skypolaris.org/wp-content/uploads/IGS%20Files/)(.*?)(%20)(.*?)(.igc)$", u)
	if check == true {
		return true
	}
	return false
}

// The following function compares two times and is taken entirely from stackoverflow.
// https://stackoverflow.com/questions/36530251/golang-time-since-with-months-and-years/36531443#36531443
func diff(a, b time.Time) (year, month, day, hour, min, sec int) {
	if a.Location() != b.Location() {
		b = b.In(a.Location())
	}
	if a.After(b) {
		a, b = b, a
	}
	y1, M1, d1 := a.Date()
	y2, M2, d2 := b.Date()
	h1, m1, s1 := a.Clock()
	h2, m2, s2 := b.Clock()
	year = int(y2 - y1)
	month = int(M2 - M1)
	day = int(d2 - d1)
	hour = int(h2 - h1)
	min = int(m2 - m1)
	sec = int(s2 - s1)
	// Normalize negative values
	if sec < 0 {
		sec += 60
		min--
	}
	if min < 0 {
		min += 60
		hour--
	}
	if hour < 0 {
		hour += 24
		day--
	}
	if day < 0 {
		// days in month:
		t := time.Date(y1, M1, 32, 0, 0, 0, 0, time.UTC)
		day += 32 - t.Day()
		month--
	}
	if month < 0 {
		month += 12
		year--
	}
	return
}
