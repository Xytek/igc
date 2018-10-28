package main

import (
	"fmt"
	"net/http"
)

// adminCount counts all tracks in the DB
func adminCount(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%v\n", db.countTrack())
}

// adminDelete deletes all tracks in the DB
func adminDelete(w http.ResponseWriter, r *http.Request) {
	count := db.countTrack()

	db.deleteAllTracks()
	fmt.Fprintf(w, "%v\n", count)
}
