package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_apiRedirect(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(apiRedirect))
	defer ts.Close()

	resp, _ := http.Get(ts.URL + "/paragliding/")
	if resp.StatusCode != http.StatusSeeOther {
		t.Errorf("Expected StatusCode %d, received %d", http.StatusSeeOther, resp.StatusCode)
		return
	}

}

func Test_apiGet(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(apiGet))
	defer ts.Close()

	meta := Meta{"P0Y0M0DT0H0M0S", "Service for Paragliding tracks.", "v1"}
	resp, err := http.Get(ts.URL + "/paragliding/api/track/{igcId}")
	if err != nil {
		t.Errorf("Error making the GET request, %s", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected StatusCode %d, received %d", http.StatusOK, resp.StatusCode)
		return
	}
	var a Meta
	err = json.NewDecoder(resp.Body).Decode(&a)
	if err != nil {
		t.Errorf("Error parsing the expected JSON body. Got %s", err)
	}
	if a.Version != meta.Version || a.Info != meta.Info {
		t.Errorf("Expected %v got %v", meta, a)
	}
}
