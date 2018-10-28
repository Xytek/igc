package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/gorilla/mux"
)

// Webhook struct storing users for the database
type Webhook struct {
	URL             string        `json:"webhookURL"`
	MinTriggerValue int           `json:"minTriggerValue"`
	ID              bson.ObjectId `json:"-" bson:"_id"`
	LastCheck       int64         `json:"-" bson:"lastCheck"`
}

// WebhookMsg struct for posting messages to users
type WebhookMsg struct {
	Content    string        `json:"content"`
	TLatest    int64         `json:"t_latest"`
	Tracks     []int         `json:"tracks"`
	Processing time.Duration `json:"processing"`
}

// webhookNewTrackGet
func webhookNewTrackPost(w http.ResponseWriter, r *http.Request) {
	// Stores the user data that we're decoding
	var user Webhook

	// Defaults MinTriggerValue to 1 if it's not changed in the decode
	user.MinTriggerValue = 1

	// Read in the user data
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Couldn't decode the body", http.StatusBadRequest)
		return
	}

	// Ensure the user has filled something as URL
	if user.URL == "" {
		http.Error(w, "The request does not fulfill the requirements", http.StatusBadRequest)
		return
	}

	// Add the user to the database
	user = db.addWebhook(user)

	// Print the id of the new user
	fmt.Fprintf(w, "%s", user.ID.Hex())
}

// webhookNewTrackIDGet
func webhookNewTrackIDGet(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["webhook_id"]

	wh := db.getWebhook(id)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(wh)
}

// webhookNewTrackIDDelete
func webhookNewTrackIDDelete(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["webhook_id"]

	// Save the webhook temporarily
	wh := db.getWebhook(id)

	// Delete the webhook from the database
	db.deleteWebhook(id)

	// Return the information about the now deleted webhook
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(wh)

}

// CheckUsers loops through each user and checks if it's time to notify them
func CheckUsers() {
	funcStart := time.Now()

	// Get all webhooks
	webhooks := db.getAllWebhooks()
	var whMsg WebhookMsg
	// Cycle through all tracks and fetch tracks with timestamp newer than the webhooks (edigble tracks)
	for _, u := range webhooks {
		tracks := db.getTracksAfter(u.LastCheck)
		for _, t := range tracks {
			whMsg.Tracks = append(whMsg.Tracks, t.SimpleID)
		}
		// If the length of the new tracks matches the users min trigger value then he'll be notified
		if len(tracks)%u.MinTriggerValue == 0 {
			// Notify the user
			NotifyUser(u, whMsg, u.URL, funcStart)
		}
	}
}

// NotifyUser sends a notification to the users
func NotifyUser(u Webhook, whMsg WebhookMsg, url string, funcStart time.Time) {
	// Find the latest added tracks timestamp
	whMsg.TLatest = db.latestTicker()

	// Updates the latest timestamp of the webhook in the database
	db.updateWebhookTimestamp(u, whMsg.TLatest)

	// Creates a string in the format 1, 2, 3 from the id array
	tracksArrayPresentable := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(whMsg.Tracks)), ", "), "[]")

	// Checks the processing time since the start of the function (FindUsers) in ms
	whMsg.Processing = time.Since(funcStart) / 1000000

	// Creates the string that'll be shown to the user
	whMsg.Content = fmt.Sprintf("Latest timestamp: %v, %v new tracks are: %s. (processing: %vs %vms)",
		whMsg.TLatest, len(whMsg.Tracks), tracksArrayPresentable, int(whMsg.Processing/1000), int(whMsg.Processing%1000))

	msg, _ := json.Marshal(whMsg)

	// POST the message
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(msg))
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
}
