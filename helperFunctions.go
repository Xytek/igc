package main

import (
	"regexp"
	"time"
)

// IDs is an array holding the internal ID system
var IDs []int

// lastUsed is an int that tracks the ID of the last added track
var lastUsed int

// Used for calculating the applications uptime
var startTime = time.Now()

// IDsInit adds each existing ID from the database in the IDs array
func IDsInit() {
	// Gets all the tracks
	var t = db.getAllTracks()

	// Adds each existing ID in the IDs array
	for _, track := range t {
		IDs = append(IDs, track.SimpleID)
		lastUsed = track.SimpleID + 1
	}
}

// checkURL uses a regular expression to check that the URL is of the type that this program can read.
func checkURL(u string) bool {
	// Had a larger test initially as I could only find two different types of URLs to .igc files, but decided on a simpler approach that only checks that it ends with .igc
	check, _ := regexp.MatchString("(.igc)$", u)
	// Below is the old one
	// There are two types of files that can be read. Only difference is one uses %20to%20 and other uses %20 between the cities. Therefore the rest is hardcoded
	// check, _ := regexp.MatchString("^(http://skypolaris.org/wp-content/uploads/IGS%20Files/)(.*?)(%20)(.*?)(.igc)$", u)
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
