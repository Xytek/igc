package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

func Test_webhookNewTrackPost(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(webhookNewTrackPost))
	defer ts.Close()
	DBInitTest()

	// Test proper body, actual data input and response
	test := "{\"webhookURL\":\"https://discordapp.com/api/webhooks/506087941764939790/4CuwPYPKQy1UylbdbOuf2fgcoL5Dmtk7-W7D3yNrQQ9P_bxxKUL0dT5AjrNI81i_b8yd\"}"
	resp, err := http.Post(ts.URL+"/paragliding/api/webhook/new_track", "application/json", strings.NewReader(test))
	if err != nil {
		t.Errorf("Error creating the POST request, %s", err)
	}
	if resp.StatusCode != http.StatusOK {
		all, _ := ioutil.ReadAll(resp.Body)
		t.Errorf("Expected StatusCode %d, received %d, Body: %s", http.StatusOK, resp.StatusCode, all)
	}

	a, _ := ioutil.ReadAll(resp.Body)

	if err != nil {
		t.Errorf("Error parsing the expected JSON body. Got %s", err)
	}

	wh := db.getWebhook(string(a))
	print(wh.MinTriggerValue)
	if wh.MinTriggerValue != 1 || wh.URL != "https://discordapp.com/api/webhooks/506087941764939790/4CuwPYPKQy1UylbdbOuf2fgcoL5Dmtk7-W7D3yNrQQ9P_bxxKUL0dT5AjrNI81i_b8yd" {
		t.Error("Something went wrong when adding the webhook to the database")
	}
	db.deleteWebhook(string(a))
}

func Test_webhookNewTrackIDGet(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/paragliding/api/webhook/new_track/{webhook_id}", webhookNewTrackIDGet).Methods("GET")
	ts := httptest.NewServer(r)
	defer ts.Close()
	DBInitTest()

	var tData Webhook
	tData.URL = "https://discordapp.com/api/webhooks/506087941764939790/4CuwPYPKQy1UylbdbOuf2fgcoL5Dmtk7-W7D3yNrQQ9P_bxxKUL0dT5AjrNI81i_b8yd"
	tData = db.addWebhook(tData)

	resp, err := http.Get(ts.URL + "/paragliding/api/webhook/new_track/" + tData.ID.Hex())
	if err != nil {
		t.Errorf("Error making the GET request, %s", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected StatusCode %d, received %d", http.StatusOK, resp.StatusCode)
		return
	}
	var a Webhook
	err = json.NewDecoder(resp.Body).Decode(&a)
	if a.URL != tData.URL {
		t.Errorf("Expected %v, got %v", tData, a)
	}
	db.deleteWebhook(tData.ID.Hex())
}

// func Test_webhookNewTrackIDDelete(t *testing.T) {
// 	r := mux.NewRouter()
// 	r.HandleFunc("/paragliding/api/webhook/new_track/{webhook_id}", webhookNewTrackIDDelete).Methods("DELETE")
// 	ts := httptest.NewServer(r)
// 	defer ts.Close()
// 	DBInitTest()

// 	var tData Webhook
// 	tData.URL = "https://discordapp.com/api/webhooks/506087941764939790/4CuwPYPKQy1UylbdbOuf2fgcoL5Dmtk7-W7D3yNrQQ9P_bxxKUL0dT5AjrNI81i_b8yd"
// 	tData = db.addWebhook(tData)
// 	count1 := len(db.getAllWebhooks())
// 	//resp, err := http.Post(ts.URL+"/paragliding/api/webhook/new_track", "application/json", strings.NewReader(test))
// 	resp, err := http.Delete(ts.URL+"/paragliding/api/webhook/new_track/"+tData.ID.Hex()), "application/json", strings.NewReader()
// 	if err != nil {
// 		t.Errorf("Error making the DELETE request, %s", err)
// 	}
// 	if resp.StatusCode != http.StatusOK {
// 		t.Errorf("Expected StatusCode %d, received %d", http.StatusOK, resp.StatusCode)
// 		return
// 	}
// 	var a Webhook
// 	err = json.NewDecoder(resp.Body).Decode(&a)
// 	if a.URL != tData.URL {
// 		t.Errorf("Expected %v, got %v", tData, a)
// 	}
// 	count2 := len(db.getAllWebhooks())

// 	if count1 <= count2 {
// 		t.Errorf("Expected %d webhooks after delete, but have %d", count2, count1)
// 	}
// }
