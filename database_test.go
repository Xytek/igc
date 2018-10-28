package main

import (
	"testing"

	igc "github.com/marni/goigc"
)

/*
The following function tests all of these functions that does something with track items in the database
addTrack
countTrack
getTrack
getAllTracks
getTracksAfter
latestTicker
deleteAllTracks
*/

func Test_TrackDatabase(t *testing.T) {
	// Make sure the test runs the test database
	DBInitTest()
	// Make sure IDs and db don't carry data from other tests
	IDs = nil
	db.deleteAllTracks()
	// Check that we can get data from URL
	url := "http://skypolaris.org/wp-content/uploads/IGS%20Files/Madrid%20to%20Jerez.igc"
	track, _ := igc.ParseLocation(url)

	// Add testing data
	tData := Track{"A date", "a pilot", "type", "GliderID", 0, url, 10, 0}
	db.addTrack(tData)
	// Add one more track
	tData = Track{"A date", track.Pilot, "type", "GliderID", 0, url, 0, 0}
	tData = db.addTrack(tData)

	if db.countTrack() != 2 {
		t.Error("Problem with addTrack(tData) or countTrack()")
	}

	if len(IDs) != 2 || IDs[1] != 1 {
		t.Error("Problem with IDs[]")
	}

	tr, _ := db.getTrack(1)
	if tr.Pilot != tData.Pilot {
		t.Error("Problem with getTrack(0)")
	}

	tracks := db.getAllTracks()

	if len(tracks) != 2 || tracks[1].Pilot != tData.Pilot {
		t.Error("Problem with getAllTracks(0)")
	}

	tracks = db.getTracksAfter(20)
	if len(tracks) != 1 || tracks[0].Pilot != track.Pilot {
		t.Error("Problem with getTracksAfter(20)")
	}

	ts := db.latestTicker()
	if ts != tData.Timestamp {
		t.Errorf("Problem with latestTicker(): Expected %v but got %v", tData.Timestamp, ts)
	}

	db.deleteAllTracks()
	tracks = db.getAllTracks()

	if len(tracks) != 0 {
		t.Error("Problem with deleteAllTracks()")
	}
}

/*
The following function tests all of these functions that does something with webhook items in the database
addWebhook
getWebhook
deleteWebhook
getAllWebhooks
updateWebhookTimestamp
*/
func Test_WebhookDatabase(t *testing.T) {
	// Make sure the test runs the test database
	DBInitTest()
	// Make sure IDs don't carry data from other tests
	IDs = nil
	url := "https://discordapp.com/api/webhooks/506087941764939790/4CuwPYPKQy1UylbdbOuf2fgcoL5Dmtk7-W7D3yNrQQ9P_bxxKUL0dT5AjrNI81i_b8yd"
	// Add testing data
	var tData Webhook
	tData.URL = url
	tData.MinTriggerValue = 1

	wh := db.addWebhook(tData)
	if wh.URL != tData.URL || wh.MinTriggerValue != 1 {
		t.Error("Problem with addWebhook(tData)")
	}

	wh2 := db.getWebhook(wh.ID.Hex())
	if wh2 != wh {
		t.Error("Problem with getWebhook(id)")
	}
	wh2 = db.addWebhook(wh2)

	whs := db.getAllWebhooks()
	if len(whs) != 2 {
		t.Error("Problem with getAllWebhooks()")
	}

	db.updateWebhookTimestamp(wh2, 10)
	wh2 = db.getWebhook(wh2.ID.Hex())

	if wh.LastCheck == wh2.LastCheck || wh2.LastCheck != 10 {
		t.Error("Problem with updateWebhookTimestamp()")
	}

	db.deleteWebhook(wh2.ID.Hex())
	db.deleteWebhook(wh.ID.Hex())
	whs = db.getAllWebhooks()
	if len(whs) != 0 {
		t.Error("Problem with deleteWebhook()")
	}
}
