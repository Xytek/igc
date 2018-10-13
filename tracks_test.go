package main

import (
	"testing"

	igc "github.com/marni/goigc"
)

func Test_addTrack(t *testing.T) {
	db := &TrackDB{}
	track, _ := igc.ParseLocation("http://skypolaris.org/wp-content/uploads/IGS%20Files/Madrid%20to%20Jerez.igc")
	tData := Track{track.Header.Date, track.Pilot, track.GliderType, track.GliderID, track.Points[0].Distance(track.Points[len(track.Points)-1])}
	db.Init()
	db.Add(tData, ID{"ID0"})
	if db.Count() != 1 {
		t.Error("Wrong track count")
	}
	if len(IDs) != 1 {
		t.Error("Wrong IDs count")
	}
	tr, _ := db.Get("ID0")
	if tr.Pilot != tData.Pilot {
		t.Error("The track was not added.")
	}
}
